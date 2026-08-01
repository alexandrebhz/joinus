#!/usr/bin/env python3
"""
Scrape open jobs from JoinUs startup websites and create postings in the internal DB.

Flow:
  1. List startups from JoinUs API (same data as https://joinus.ie/startups)
  2. For each startup, follow its website (Visit Website)
  3. Use Firecrawl map + agent to find and extract open job listings
  4. POST new jobs to JoinUs API (POST /api/v1/jobs)

Pilot usage (first 10 startups, dry-run):
  python scrape_startup_jobs.py --limit 10 --dry-run

Create jobs (requires auth):
  export JOINUS_API_BASE='https://joinus-production.up.railway.app'
  export JOINUS_ACCESS_TOKEN='<jwt>'   # or JOINUS_EMAIL + JOINUS_PASSWORD
  python scrape_startup_jobs.py --limit 10

Requires: firecrawl-cli (npm i -g firecrawl-cli), httpx
"""

from __future__ import annotations

import argparse
import json
import os
import re
import subprocess
import sys
import time
from pathlib import Path
from typing import Any, Optional
from urllib.parse import urlparse

import httpx

DEFAULT_API_BASE = "https://joinus-production.up.railway.app"
USER_AGENT = "joinus-startup-job-scraper/1.0"
SCRIPT_DIR = Path(__file__).resolve().parent
FIRECRAWL_DIR = SCRIPT_DIR.parent / ".firecrawl"
SCHEMA_PATH = SCRIPT_DIR / "job_schema.json"
STATE_PATH = SCRIPT_DIR / "startup_jobs_state.json"

JOB_TYPE_MAP = {
    "full-time": "full_time",
    "fulltime": "full_time",
    "full_time": "full_time",
    "part-time": "part_time",
    "parttime": "part_time",
    "part_time": "part_time",
    "contract": "contract",
    "contractor": "contract",
    "freelance": "contract",
    "intern": "internship",
    "internship": "internship",
}

LOCATION_TYPE_MAP = {
    "remote": "remote",
    "hybrid": "hybrid",
    "onsite": "onsite",
    "on-site": "onsite",
    "on site": "onsite",
    "office": "onsite",
    "in-office": "onsite",
}

COUNTRY_CURRENCY = {
    "united states": "USD",
    "usa": "USD",
    "us": "USD",
    "united kingdom": "GBP",
    "uk": "GBP",
    "ireland": "EUR",
    "germany": "EUR",
    "france": "EUR",
    "romania": "RON",
    "canada": "CAD",
    "australia": "AUD",
}


def normalize_api_origin(raw: str) -> str:
    s = (raw or "").strip().rstrip("/")
    if not s:
        return ""
    if "://" not in s:
        s = "https://" + s.lstrip("/")
    p = urlparse(s)
    if not p.netloc:
        return ""
    return f"{(p.scheme or 'https').lower()}://{p.netloc}"


def slugify(name: str) -> str:
    s = name.lower().strip()
    s = re.sub(r"[^a-z0-9]+", "-", s)
    return s.strip("-") or "startup"


def ensure_firecrawl_dir() -> Path:
    FIRECRAWL_DIR.mkdir(parents=True, exist_ok=True)
    return FIRECRAWL_DIR


CAREER_URL_KEYWORDS = (
    "career",
    "job",
    "hiring",
    "opening",
    "join-us",
    "joinus",
    "work-with-us",
    "greenhouse.io",
    "lever.co",
    "ashbyhq.com",
    "workable.com",
    "breezy.hr",
    "recruitee.com",
)

CAREER_PATH_SUFFIXES = (
    "/careers",
    "/jobs",
    "/company/careers",
    "/about/careers",
    "/work-with-us",
    "/join-us",
)


EXCLUDED_URL_PARTS = (
    "/blog/",
    "/insights/",
    "/leads/",
    "how-to-",
    "opening-lines",
    "workflow",
    "beta-sign-up",
    "job-posting-data",
    "recruiting-deep-dive",
)


def filter_career_urls(urls: list[str], website: str) -> list[str]:
    scored: list[tuple[int, str]] = []
    for url in urls:
        low = url.lower()
        if any(part in low for part in EXCLUDED_URL_PARTS):
            continue
        score = 0
        if low.rstrip("/").endswith("/careers") or low.rstrip("/").endswith("/jobs"):
            score += 5
        if "ashbyhq.com" in low or "greenhouse.io" in low or "lever.co" in low:
            score += 6
        score += sum(1 for kw in CAREER_URL_KEYWORDS if kw in low)
        if score:
            scored.append((score, url))
    scored.sort(key=lambda x: (-x[0], x[1]))
    ordered = [u for _, u in scored]
    if ordered:
        return ordered
    base = website.rstrip("/")
    return [base + suffix for suffix in CAREER_PATH_SUFFIXES]


def run_firecrawl(args: list[str], *, timeout: int = 300) -> subprocess.CompletedProcess[str]:
    cmd = ["firecrawl", *args]
    return subprocess.run(
        cmd,
        capture_output=True,
        text=True,
        timeout=timeout,
        check=False,
    )


def strip_ansi(text: str) -> str:
    return re.sub(r"\x1b\[[0-9;]*[A-Za-z]", "", text or "").strip()


def firecrawl_ready() -> bool:
    proc = run_firecrawl(["--version", "--auth-status"], timeout=30)
    out = (proc.stdout or "") + (proc.stderr or "")
    return "authenticated: true" in out.lower()


def load_state(path: Path) -> dict[str, Any]:
    if not path.is_file():
        return {"processed_startups": {}, "created_jobs": {}}
    try:
        with path.open(encoding="utf-8") as f:
            data = json.load(f)
        if isinstance(data, dict):
            data.setdefault("processed_startups", {})
            data.setdefault("created_jobs", {})
            return data
    except (OSError, json.JSONDecodeError):
        pass
    return {"processed_startups": {}, "created_jobs": {}}


def save_state(path: Path, state: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8") as f:
        json.dump(state, f, indent=2, sort_keys=True)


def request_with_retry(
    client: httpx.Client,
    method: str,
    url: str,
    *,
    attempts: int = 6,
    backoff: float = 1.5,
    **kwargs: Any,
) -> httpx.Response:
    last: Optional[BaseException] = None
    for i in range(max(1, attempts)):
        try:
            return client.request(method, url, **kwargs)
        except (httpx.ConnectError, httpx.TimeoutException) as exc:
            last = exc
            if i + 1 >= attempts:
                break
            wait = min(backoff * (2**i), 45.0)
            print(f"warn: {method} {url!r} {exc}; retry in {wait:.1f}s", flush=True)
            time.sleep(wait)
    assert last is not None
    raise last


def api_json(
    client: httpx.Client,
    method: str,
    url: str,
    *,
    token: str,
    json_body: Any = None,
) -> dict[str, Any]:
    headers = {"Authorization": f"Bearer {token}", "User-Agent": USER_AGENT}
    kw: dict[str, Any] = {"headers": headers}
    if json_body is not None:
        kw["json"] = json_body
    r = request_with_retry(client, method, url, **kw)
    try:
        r.raise_for_status()
    except httpx.HTTPStatusError as exc:
        detail = ""
        try:
            detail = exc.response.text[:4000]
        except Exception:
            pass
        raise RuntimeError(f"{method} {url} -> HTTP {exc.response.status_code}: {detail}") from exc
    body = r.json()
    if not body.get("success"):
        raise RuntimeError(f"{method} {url} -> {body!r}")
    data = body.get("data")
    if data is None:
        raise RuntimeError(f"{method} {url} missing data")
    return data if isinstance(data, dict) else {"items": data}


def login(client: httpx.Client, base: str, email: str, password: str) -> str:
    r = request_with_retry(
        client,
        "POST",
        f"{base}/api/v1/auth/login",
        json={"email": email, "password": password},
    )
    r.raise_for_status()
    body = r.json()
    if not body.get("success"):
        raise RuntimeError(f"login failed: {body!r}")
    token = (body.get("data") or {}).get("access_token")
    if not token:
        raise RuntimeError("login response missing access_token")
    return str(token)


def fetch_startups(client: httpx.Client, base: str, *, limit: int, page: int) -> list[dict[str, Any]]:
    page_size = min(limit, 100) if limit > 0 else 12
    url = f"{base}/api/v1/startups?page={page}&page_size={page_size}"
    r = request_with_retry(client, "GET", url)
    r.raise_for_status()
    body = r.json()
    if not body.get("success"):
        raise RuntimeError(f"list startups failed: {body!r}")
    startups = body.get("data") or []
    if limit > 0:
        return startups[:limit]
    return startups


def fetch_existing_jobs(client: httpx.Client, base: str, startup_id: str) -> list[dict[str, Any]]:
    url = f"{base}/api/v1/jobs?startup_id={startup_id}&page=1&page_size=100"
    r = request_with_retry(client, "GET", url)
    r.raise_for_status()
    body = r.json()
    return body.get("data") or []


def discover_careers_urls(website: str, slug: str, out_dir: Path) -> list[str]:
    """Use Firecrawl map to find likely careers/jobs pages on the company site."""
    out_file = out_dir / f"map-{slug}.txt"
    proc = run_firecrawl(
        [
            "map",
            website,
            "--search",
            "careers jobs hiring openings",
            "--limit",
            "40",
            "-o",
            str(out_file),
        ],
        timeout=120,
    )
    urls: list[str] = []
    if out_file.is_file():
        for line in out_file.read_text(encoding="utf-8").splitlines():
            u = line.strip()
            if u.startswith("http"):
                urls.append(u)
    elif proc.returncode != 0:
        print(
            f"warn: map failed for {website}: {strip_ansi(proc.stderr or proc.stdout)}",
            flush=True,
        )

    return filter_career_urls(urls, website)[:6]


def parse_agent_output(raw: Any) -> dict[str, Any]:
    if isinstance(raw, dict):
        # Firecrawl envelope: { success, status, data: { has_open_jobs, jobs } }
        if raw.get("success") and isinstance(raw.get("data"), dict):
            inner = raw["data"]
            if "jobs" in inner or "has_open_jobs" in inner:
                return inner
        for key in ("data", "result", "extracted", "output"):
            val = raw.get(key)
            if isinstance(val, dict) and ("jobs" in val or "has_open_jobs" in val):
                return val
        if "jobs" in raw or "has_open_jobs" in raw:
            return raw
    raise RuntimeError(f"unexpected agent output: {raw!r}")


def extract_jobs_with_agent(
    *,
    company_name: str,
    website: str,
    careers_urls: list[str],
    slug: str,
    out_dir: Path,
    max_credits: int,
) -> dict[str, Any]:
    out_file = out_dir / f"agent-{slug}.json"
    primary = careers_urls[0] if careers_urls else website
    prompt = (
        f"On {company_name}'s website, find all currently open job positions. "
        f"Start from {primary} and follow links to external ATS boards (Greenhouse, Lever, Ashby) if needed. "
        "Extract each role's title, full description, requirements, location, job type, and application URL/email. "
        "Only include jobs that appear actively open. If none, return has_open_jobs=false and jobs=[]."
    )
    proc = run_firecrawl(
        [
            "agent",
            prompt,
            "--urls",
            primary,
            "--schema-file",
            str(SCHEMA_PATH),
            "--model",
            "spark-1-mini",
            "--max-credits",
            str(max_credits),
            "--wait",
            "--timeout",
            "180",
            "-o",
            str(out_file),
            "--json",
        ],
        timeout=240,
    )

    if out_file.is_file():
        try:
            return parse_agent_output(json.loads(out_file.read_text(encoding="utf-8")))
        except (json.JSONDecodeError, RuntimeError) as exc:
            if proc.returncode != 0:
                raise RuntimeError(
                    f"firecrawl agent failed for {company_name}: {strip_ansi(proc.stderr or proc.stdout)}"
                ) from exc
            raise

    if proc.returncode != 0:
        raise RuntimeError(
            f"firecrawl agent failed for {company_name}: {strip_ansi(proc.stderr or proc.stdout)}"
        )
    raise RuntimeError(f"firecrawl agent produced no output for {company_name}")


def normalize_job_type(raw: str) -> str:
    key = (raw or "").strip().lower()
    return JOB_TYPE_MAP.get(key, "full_time")


def normalize_location_type(raw: str) -> str:
    key = (raw or "").strip().lower()
    return LOCATION_TYPE_MAP.get(key, "remote")


def default_currency(location: str) -> str:
    loc = (location or "").strip().lower()
    for key, cur in COUNTRY_CURRENCY.items():
        if key in loc:
            return cur
    return "USD"


def country_from_startup(startup: dict[str, Any], job_country: str) -> str:
    c = (job_country or "").strip()
    if c:
        return c
    loc = (startup.get("location") or "").strip()
    return loc or "Unknown"


def pad_text(text: str, minimum: int, fallback: str) -> str:
    t = (text or "").strip()
    if len(t) >= minimum:
        return t
    if t:
        return t + (" " * (minimum - len(t)))
    return fallback


def build_job_payload(startup: dict[str, Any], job: dict[str, Any]) -> dict[str, Any]:
    title = pad_text(job.get("title") or "", 5, "Open Role")
    description = pad_text(job.get("description") or "", 20, "Role details to be updated.")
    requirements = pad_text(
        job.get("requirements") or "",
        10,
        "See full job description for requirements.",
    )

    country = country_from_startup(startup, job.get("country") or "")
    currency = (job.get("currency") or "").strip() or default_currency(country)

    payload: dict[str, Any] = {
        "startup_id": startup["id"],
        "title": title[:100],
        "description": description,
        "requirements": requirements,
        "job_type": normalize_job_type(job.get("job_type") or ""),
        "location_type": normalize_location_type(job.get("location_type") or ""),
        "city": (job.get("city") or "").strip(),
        "country": country,
        "currency": currency[:3] if currency else "USD",
    }

    if job.get("salary_min") is not None:
        payload["salary_min"] = int(job["salary_min"])
    if job.get("salary_max") is not None:
        payload["salary_max"] = int(job["salary_max"])

    app_url = (job.get("application_url") or "").strip()
    if app_url.startswith("http"):
        payload["application_url"] = app_url

    app_email = (job.get("application_email") or "").strip()
    if "@" in app_email:
        payload["application_email"] = app_email

    return payload


def dedupe_key(startup_id: str, job: dict[str, Any]) -> str:
    title = (job.get("title") or "").strip().lower()
    app = (job.get("application_url") or job.get("application_email") or "").strip().lower()
    return f"{startup_id}:{title}:{app}"


def process_startup(
    *,
    client: httpx.Client,
    api_base: str,
    token: Optional[str],
    startup: dict[str, Any],
    state: dict[str, Any],
    out_dir: Path,
    dry_run: bool,
    skip_existing: bool,
    max_credits: int,
    delay: float,
) -> dict[str, Any]:
    sid = startup["id"]
    name = startup.get("name") or "Unknown"
    website = (startup.get("website") or "").strip()
    slug = startup.get("slug") or slugify(name)

    result = {
        "startup_id": sid,
        "startup_name": name,
        "website": website,
        "careers_urls": [],
        "jobs_found": 0,
        "jobs_created": 0,
        "jobs_skipped": 0,
        "errors": [],
    }

    if not website:
        result["errors"].append("no website")
        return result

    if skip_existing and sid in state.get("processed_startups", {}):
        cached = state["processed_startups"][sid]
        print(f"skip cached {name!r} ({cached.get('jobs_found', 0)} jobs previously)", flush=True)
        return cached

    print(f"\n=== {name} ===", flush=True)
    print(f"website: {website}", flush=True)

    careers_urls = discover_careers_urls(website, slug, out_dir)
    result["careers_urls"] = careers_urls
    print(f"careers candidates: {len(careers_urls)}", flush=True)
    for u in careers_urls[:3]:
        print(f"  - {u}", flush=True)

    try:
        extracted = extract_jobs_with_agent(
            company_name=name,
            website=website,
            careers_urls=careers_urls,
            slug=slug,
            out_dir=out_dir,
            max_credits=max_credits,
        )
    except Exception as exc:
        msg = strip_ansi(str(exc))
        if len(msg) > 500:
            msg = msg[:500] + "..."
        result["errors"].append(msg)
        print(f"error extracting jobs: {msg}", flush=True)
        state.setdefault("processed_startups", {})[sid] = result
        return result

    jobs = extracted.get("jobs") or []
    has_open = extracted.get("has_open_jobs", bool(jobs))
    result["jobs_found"] = len(jobs)
    print(f"open jobs found: {len(jobs)} (has_open_jobs={has_open})", flush=True)

    if dry_run or not jobs:
        for j in jobs[:5]:
            print(f"  [dry] {j.get('title')}", flush=True)
        state.setdefault("processed_startups", {})[sid] = result
        return result

    if not token:
        result["errors"].append("missing auth token for job creation")
        return result

    existing = fetch_existing_jobs(client, api_base, sid)
    existing_titles = {(j.get("title") or "").strip().lower() for j in existing}

    for job in jobs:
        key = dedupe_key(sid, job)
        if key in state.get("created_jobs", {}):
            result["jobs_skipped"] += 1
            continue

        title_l = (job.get("title") or "").strip().lower()
        if title_l in existing_titles:
            result["jobs_skipped"] += 1
            print(f"  skip existing: {job.get('title')}", flush=True)
            continue

        payload = build_job_payload(startup, job)
        try:
            created = api_json(
                client,
                "POST",
                f"{api_base}/api/v1/jobs",
                token=token,
                json_body=payload,
            )
            job_id = created.get("id")
            state.setdefault("created_jobs", {})[key] = job_id
            result["jobs_created"] += 1
            print(f"  created: {payload['title']} -> {job_id}", flush=True)
        except Exception as exc:
            result["errors"].append(f"{job.get('title')}: {exc}")
            print(f"  failed: {job.get('title')}: {exc}", flush=True)

        if delay > 0:
            time.sleep(delay)

    state.setdefault("processed_startups", {})[sid] = result
    return result


def main() -> int:
    ap = argparse.ArgumentParser(
        description="Scrape startup websites for open jobs and create JoinUs postings."
    )
    ap.add_argument("--api-base", default=os.environ.get("JOINUS_API_BASE", DEFAULT_API_BASE))
    ap.add_argument("--token", default=os.environ.get("JOINUS_ACCESS_TOKEN", ""))
    ap.add_argument("--email", default=os.environ.get("JOINUS_EMAIL", ""))
    ap.add_argument("--password", default=os.environ.get("JOINUS_PASSWORD", ""))
    ap.add_argument("--limit", type=int, default=10, help="Max startups to process (default: 10)")
    ap.add_argument("--page", type=int, default=1, help="Startups list page (default: 1)")
    ap.add_argument("--dry-run", action="store_true", help="Scrape only; do not POST jobs")
    ap.add_argument("--skip-existing", action="store_true", help="Skip startups already in state file")
    ap.add_argument("--state", type=Path, default=STATE_PATH)
    ap.add_argument("--max-credits", type=int, default=150, help="Firecrawl agent max credits per startup")
    ap.add_argument("--delay", type=float, default=0.5, help="Delay between job POSTs")
    ap.add_argument("--startup-delay", type=float, default=2.0, help="Delay between startups")
    args = ap.parse_args()

    api_base = normalize_api_origin(args.api_base)
    if not api_base:
        print("error: set --api-base or JOINUS_API_BASE", file=sys.stderr)
        return 1

    out_dir = ensure_firecrawl_dir() / "startup-jobs-pilot"
    out_dir.mkdir(parents=True, exist_ok=True)

    if not firecrawl_ready():
        print(
            "error: firecrawl is not authenticated. Run: firecrawl login --browser",
            file=sys.stderr,
        )
        return 1

    state = load_state(args.state.resolve())
    token = (args.token or "").strip()

    summary: list[dict[str, Any]] = []

    with httpx.Client(timeout=90.0, headers={"User-Agent": USER_AGENT}) as client:
        startups = fetch_startups(client, api_base, limit=args.limit, page=args.page)
        print(f"Fetched {len(startups)} startups from {api_base}", flush=True)

        if not args.dry_run and not token:
            if args.email and args.password:
                token = login(client, api_base, args.email, args.password)
            else:
                print(
                    "error: for job creation set JOINUS_ACCESS_TOKEN or JOINUS_EMAIL+JOINUS_PASSWORD",
                    file=sys.stderr,
                )
                return 1

        for i, startup in enumerate(startups):
            result = process_startup(
                client=client,
                api_base=api_base,
                token=token or None,
                startup=startup,
                state=state,
                out_dir=out_dir,
                dry_run=args.dry_run,
                skip_existing=args.skip_existing,
                max_credits=args.max_credits,
                delay=args.delay,
            )
            summary.append(result)
            save_state(args.state.resolve(), state)
            if args.startup_delay > 0 and i + 1 < len(startups):
                time.sleep(args.startup_delay)

    report_path = out_dir / "pilot-report.json"
    with report_path.open("w", encoding="utf-8") as f:
        json.dump(summary, f, indent=2)

    total_found = sum(r.get("jobs_found", 0) for r in summary)
    total_created = sum(r.get("jobs_created", 0) for r in summary)
    total_skipped = sum(r.get("jobs_skipped", 0) for r in summary)
    with_errors = [r for r in summary if r.get("errors")]

    print("\n=== PILOT SUMMARY ===", flush=True)
    print(f"startups processed: {len(summary)}", flush=True)
    print(f"jobs found:         {total_found}", flush=True)
    print(f"jobs created:       {total_created}", flush=True)
    print(f"jobs skipped:       {total_skipped}", flush=True)
    print(f"startups w/ errors: {len(with_errors)}", flush=True)
    print(f"report: {report_path}", flush=True)
    print(f"state:  {args.state.resolve()}", flush=True)

    return 0


if __name__ == "__main__":
    raise SystemExit(main())

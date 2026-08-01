#!/usr/bin/env python3
"""Generate SQL INSERT statements for Clay jobs from Ashby API + clay.com listing."""

from __future__ import annotations

import json
import re
import subprocess
import uuid
from pathlib import Path

STARTUP_ID = "0455617e-1070-4a2e-a38d-a5ca020222e9"
ASHBY_URL = "https://api.ashbyhq.com/posting-api/job-board/claylabs"
SCRIPT_DIR = Path(__file__).resolve().parent
LISTING_PATH = SCRIPT_DIR.parent / ".firecrawl/clay-jobs/jobs-page.json"
OUT_PATH = SCRIPT_DIR / "clay_jobs_import.sql"

JOB_TYPE_MAP = {
    "fulltime": "full_time",
    "parttime": "part_time",
    "contract": "contract",
    "internship": "internship",
}


def sql_escape(s: str) -> str:
    return (s or "").replace("'", "''")


def parse_salary(text: str) -> tuple[int | None, int | None]:
    m = re.search(r"\$(\d+)K\s*[–-]\s*\$(\d+)K", text)
    if m:
        return int(m.group(1)) * 1000, int(m.group(2)) * 1000
    return None, None


def extract_requirements(desc: str) -> str:
    for pat in (
        r"(?i)(what you(?:'|')ll bring|who you are(?:\s+not)?|requirements|qualifications|what we(?:'|')re looking for)([\s\S]+)$",
    ):
        m = re.search(pat, desc)
        if m:
            req = m.group(0).strip()
            if len(req) >= 10:
                return req[:8000]
    parts = desc.split("\n\n")
    if len(parts) > 3:
        tail = "\n\n".join(parts[-3:])
        if len(tail) >= 10:
            return tail[:8000]
    return "See full job description for requirements."


def extract_description(desc: str, title: str) -> str:
    base_title = title.split("@")[0].strip()
    markers = (
        rf"(?i)(?:^|\n)(?:#{{1,3}}\s*)?{re.escape(base_title)}[\s\S]*",
        r"(?i)what you(?:'|')ll (?:do|achieve)[\s\S]*",
        r"(?i)about (?:the role|this role|clay)[\s\S]*",
    )
    for pat in markers:
        m = re.search(pat, desc)
        if m:
            chunk = m.group(0).strip()
            if len(chunk) >= 20:
                return chunk[:12000]
    if len(desc) >= 20:
        return desc[:12000]
    return "Role details to be updated."


def normalize_job_type(raw: str) -> str:
    key = re.sub(r"[^a-z]", "", (raw or "").lower())
    return JOB_TYPE_MAP.get(key, "full_time")


def normalize_location_type(raw: str, is_remote: bool) -> str:
    key = (raw or "").lower().replace("_", "").replace("-", "")
    if key == "onsite":
        return "onsite"
    if key == "hybrid":
        return "hybrid"
    if key == "remote" or is_remote:
        return "remote"
    return "hybrid"


def country_from_job(job: dict) -> tuple[str, str]:
    city = (job.get("location") or "").strip()
    postal = (job.get("address") or {}).get("postalAddress") or {}
    country = (postal.get("addressCountry") or "United States").strip()
    cities = [city] if city else []
    for sec in job.get("secondaryLocations") or []:
        loc = (sec.get("location") or "").strip()
        if loc and loc not in cities:
            cities.append(loc)
    return ("; ".join(cities) if cities else "")[:255], country


def salary_map_from_listing(path: Path) -> dict[str, tuple[int, int]]:
    if not path.is_file():
        return {}
    md = json.loads(path.read_text(encoding="utf-8")).get("markdown", "")
    out: dict[str, tuple[int, int]] = {}
    for m in re.finditer(r"ashby_jid=([a-f0-9-]+)", md):
        jid = m.group(1)
        if jid in out:
            continue
        chunk = md[max(0, m.start() - 350) : m.end() + 80]
        smin, smax = parse_salary(chunk)
        if smin is not None and smax is not None:
            out[jid] = (smin, smax)
    return out


def fetch_jobs() -> list[dict]:
    raw = subprocess.check_output(["curl", "-sL", ASHBY_URL], text=True)
    return json.loads(raw).get("jobs", [])


def build_insert(job: dict, salary_by_jid: dict[str, tuple[int, int]]) -> str:
    jid = job["id"]
    title = (job.get("title") or "Open Role")[:100]
    desc_plain = job.get("descriptionPlain") or ""
    description = extract_description(desc_plain, title)
    requirements = extract_requirements(desc_plain)
    if len(description) < 20:
        description = desc_plain[:12000] if desc_plain else "Role details to be updated."
    if len(requirements) < 10:
        requirements = "See full job description for requirements."

    job_type = normalize_job_type(job.get("employmentType", ""))
    location_type = normalize_location_type(job.get("workplaceType", ""), job.get("isRemote", False))
    city, country = country_from_job(job)
    smin, smax = salary_by_jid.get(jid, (None, None))
    app_url = job.get("applyUrl") or job.get("jobUrl") or f"https://www.clay.com/jobs?ashby_jid={jid}"

    cols = [
        "id", "startup_id", "title", "description", "requirements",
        "job_type", "location_type", "city", "country", "currency",
        "application_url", "status", "created_at", "updated_at",
    ]
    vals = [
        f"'{uuid.uuid4()}'",
        f"'{STARTUP_ID}'",
        f"'{sql_escape(title)}'",
        f"'{sql_escape(description)}'",
        f"'{sql_escape(requirements)}'",
        f"'{job_type}'",
        f"'{location_type}'",
        f"'{sql_escape(city)}'" if city else "NULL",
        f"'{sql_escape(country)}'",
        "'USD'",
        f"'{sql_escape(app_url)}'",
        "'active'",
        "NOW()",
        "NOW()",
    ]
    if smin is not None:
        cols.insert(9, "salary_min")
        vals.insert(9, str(smin))
    if smax is not None:
        idx = cols.index("currency")
        cols.insert(idx, "salary_max")
        vals.insert(idx, str(smax))

    lines = [
        f"-- {title}",
        f"INSERT INTO jobs ({', '.join(cols)})",
        f"VALUES ({', '.join(vals)});",
        "",
    ]
    return "\n".join(lines)


def main() -> None:
    jobs = fetch_jobs()
    salary_by_jid = salary_map_from_listing(LISTING_PATH)

    header = "\n".join([
        "-- Clay jobs import for JoinUs production",
        "-- Startup: Clay (clay-1bb6dab1)",
        f"-- startup_id: {STARTUP_ID}",
        "-- Source: https://www.clay.com/jobs (Ashby API claylabs)",
        f"-- Total jobs: {len(jobs)}",
        "",
        "BEGIN;",
        "",
    ])

    body = "".join(build_insert(j, salary_by_jid) for j in sorted(jobs, key=lambda x: x.get("title", "")))
    footer = "COMMIT;\n"

    OUT_PATH.write_text(header + body + footer, encoding="utf-8")
    with_salary = sum(1 for j in jobs if j["id"] in salary_by_jid)
    print(f"Wrote {len(jobs)} jobs to {OUT_PATH} ({with_salary} with salary)")


if __name__ == "__main__":
    main()

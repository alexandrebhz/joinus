'use client';

import { FieldRule } from '@/src/domain/entities/CrawlSite';
import { useState } from 'react';

interface FieldRuleEditorProps {
  fieldName: string;
  rule: FieldRule;
  onUpdate: (fieldName: string, rule: FieldRule) => void;
  onDelete: (fieldName: string) => void;
}

export function FieldRuleEditor({ fieldName, rule, onUpdate, onDelete }: FieldRuleEditorProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  const updateRule = (updates: Partial<FieldRule>) => {
    onUpdate(fieldName, { ...rule, ...updates });
  };

  return (
    <div className="border rounded-lg p-4 bg-gray-50">
      <div className="flex justify-between items-center mb-2">
        <div className="flex items-center gap-2">
          <h4 className="font-semibold">{fieldName}</h4>
          {rule.required && (
            <span className="text-xs bg-red-100 text-red-800 px-2 py-1 rounded">
              Required
            </span>
          )}
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setIsExpanded(!isExpanded)}
            className="text-sm text-blue-600 hover:text-blue-800"
          >
            {isExpanded ? 'Collapse' : 'Expand'}
          </button>
          <button
            type="button"
            onClick={() => onDelete(fieldName)}
            className="text-sm text-red-600 hover:text-red-800"
          >
            Delete
          </button>
        </div>
      </div>

      <div className="space-y-3">
        <div>
          <label className="block text-xs font-medium mb-1">Selector</label>
          <input
            type="text"
            value={rule.selector}
            onChange={(e) => updateRule({ selector: e.target.value })}
            className="w-full px-2 py-1 text-sm border rounded"
            placeholder=".field-selector"
          />
        </div>

        <div>
          <label className="block text-xs font-medium mb-1">Type</label>
          <select
            value={rule.type}
            onChange={(e) => updateRule({ type: e.target.value as FieldRule['type'] })}
            className="w-full px-2 py-1 text-sm border rounded"
          >
            <option value="text">Text</option>
            <option value="html">HTML</option>
            <option value="attribute">Attribute</option>
            <option value="regex">Regex</option>
          </select>
        </div>

        {isExpanded && (
          <>
            {rule.type === 'attribute' && (
              <div>
                <label className="block text-xs font-medium mb-1">Attribute Name</label>
                <input
                  type="text"
                  value={rule.attribute || ''}
                  onChange={(e) => updateRule({ attribute: e.target.value })}
                  className="w-full px-2 py-1 text-sm border rounded"
                  placeholder="href, data-id, etc."
                />
              </div>
            )}

            {rule.type === 'regex' && (
              <div>
                <label className="block text-xs font-medium mb-1">Regex Pattern</label>
                <input
                  type="text"
                  value={rule.regex_pattern || ''}
                  onChange={(e) => updateRule({ regex_pattern: e.target.value })}
                  className="w-full px-2 py-1 text-sm border rounded"
                  placeholder="\\$([0-9,]+)"
                />
              </div>
            )}

            <div>
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={rule.required}
                  onChange={(e) => updateRule({ required: e.target.checked })}
                  className="rounded"
                />
                <span className="text-xs font-medium">Required Field</span>
              </label>
            </div>

            <div>
              <label className="block text-xs font-medium mb-1">Default Value</label>
              <input
                type="text"
                value={rule.default_value || ''}
                onChange={(e) => updateRule({ default_value: e.target.value })}
                className="w-full px-2 py-1 text-sm border rounded"
                placeholder="Default if not found"
              />
            </div>

            <div>
              <label className="block text-xs font-medium mb-1">Transformations (comma-separated)</label>
              <input
                type="text"
                value={rule.transformations?.join(', ') || ''}
                onChange={(e) =>
                  updateRule({
                    transformations: e.target.value
                      .split(',')
                      .map((t) => t.trim())
                      .filter((t) => t.length > 0),
                  })
                }
                className="w-full px-2 py-1 text-sm border rounded"
                placeholder="trim, lowercase, strip_html"
              />
            </div>
          </>
        )}
      </div>
    </div>
  );
}


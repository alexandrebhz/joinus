'use client';

import { useState } from 'react';
import { FieldRule } from '@/src/domain/entities/CrawlSite';
import { FieldRuleEditor } from './FieldRuleEditor';

interface DynamicFieldsEditorProps {
  fields: Record<string, FieldRule>;
  onFieldsChange: (fields: Record<string, FieldRule>) => void;
}

export function DynamicFieldsEditor({ fields, onFieldsChange }: DynamicFieldsEditorProps) {
  const [newFieldName, setNewFieldName] = useState('');
  const [showAddField, setShowAddField] = useState(false);

  const handleAddField = () => {
    if (!newFieldName.trim()) {
      alert('Please enter a field name');
      return;
    }

    if (fields[newFieldName]) {
      alert('Field already exists');
      return;
    }

    const newField: FieldRule = {
      selector: '',
      type: 'text',
      required: false,
    };

    onFieldsChange({
      ...fields,
      [newFieldName]: newField,
    });

    setNewFieldName('');
    setShowAddField(false);
  };

  const handleUpdateField = (fieldName: string, rule: FieldRule) => {
    onFieldsChange({
      ...fields,
      [fieldName]: rule,
    });
  };

  const handleDeleteField = (fieldName: string) => {
    if (!confirm(`Are you sure you want to delete field "${fieldName}"?`)) {
      return;
    }

    const newFields = { ...fields };
    delete newFields[fieldName];
    onFieldsChange(newFields);
  };

  const handleRenameField = (oldName: string, newName: string) => {
    if (newName === oldName) return;
    if (fields[newName]) {
      alert('Field name already exists');
      return;
    }

    const newFields = { ...fields };
    newFields[newName] = newFields[oldName];
    delete newFields[oldName];
    onFieldsChange(newFields);
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold">Extraction Fields</h3>
        {!showAddField && (
          <button
            type="button"
            onClick={() => setShowAddField(true)}
            className="text-sm bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700"
          >
            + Add Field
          </button>
        )}
      </div>

      {showAddField && (
        <div className="flex gap-2 items-center p-3 bg-blue-50 rounded">
          <input
            type="text"
            value={newFieldName}
            onChange={(e) => setNewFieldName(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleAddField()}
            placeholder="Field name (e.g., company, location, salary)"
            className="flex-1 px-3 py-2 border rounded"
          />
          <button
            type="button"
            onClick={handleAddField}
            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Add
          </button>
          <button
            type="button"
            onClick={() => {
              setShowAddField(false);
              setNewFieldName('');
            }}
            className="bg-gray-200 text-gray-800 px-4 py-2 rounded hover:bg-gray-300"
          >
            Cancel
          </button>
        </div>
      )}

      <div className="space-y-3">
        {Object.entries(fields).map(([fieldName, rule]) => (
          <FieldRuleEditor
            key={fieldName}
            fieldName={fieldName}
            rule={rule}
            onUpdate={handleUpdateField}
            onDelete={handleDeleteField}
          />
        ))}

        {Object.keys(fields).length === 0 && (
          <div className="text-center py-8 text-gray-500 border-2 border-dashed rounded-lg">
            No extraction fields configured. Click "Add Field" to get started.
          </div>
        )}
      </div>
    </div>
  );
}


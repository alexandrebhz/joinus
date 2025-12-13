import { TextareaHTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string
  error?: string
  helperText?: string
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, label, error, helperText, id, ...props }, ref) => {
    const textareaId = id || `textarea-${Math.random().toString(36).substr(2, 9)}`

    return (
      <div className="w-full">
        {label && (
          <label htmlFor={textareaId} className="block text-sm font-medium text-secondary-700 mb-1">
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={textareaId}
          className={clsx(
            'w-full px-4 py-2 border rounded-lg transition-all duration-200 resize-y',
            'focus:outline-none focus:ring-2 focus:ring-offset-2',
            'disabled:bg-secondary-100 disabled:cursor-not-allowed',
            'min-h-[100px]',
            error
              ? 'border-error-500 focus:ring-error-500 focus:border-error-500'
              : 'border-secondary-300 focus:ring-primary-500 focus:border-primary-500',
            className
          )}
          {...props}
        />
        {error && <p className="mt-1 text-sm text-error-500">{error}</p>}
        {helperText && !error && <p className="mt-1 text-sm text-secondary-500">{helperText}</p>}
      </div>
    )
  }
)

Textarea.displayName = 'Textarea'

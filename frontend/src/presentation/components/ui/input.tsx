import { InputHTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  helperText?: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, label, error, helperText, id, ...props }, ref) => {
    const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`

    return (
      <div className="w-full">
        {label && (
          <label htmlFor={inputId} className="block text-sm font-medium text-secondary-700 mb-1">
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={inputId}
          className={clsx(
            'w-full px-4 py-2 border rounded-lg transition-all duration-200',
            'focus:outline-none focus:ring-2 focus:ring-offset-2',
            'disabled:bg-secondary-100 disabled:cursor-not-allowed',
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

Input.displayName = 'Input'


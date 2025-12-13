'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Button } from '@/presentation/components/ui/button'
import { Input } from '@/presentation/components/ui/input'
import { Textarea } from '@/presentation/components/ui/textarea'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/presentation/components/ui/card'
import { apiClient } from '@/infrastructure/api/api-client'
import { Mail, Send, CheckCircle2 } from 'lucide-react'

const contactSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  subject: z.string().max(500, 'Subject must be less than 500 characters').optional().or(z.literal('')),
  message: z.string().min(10, 'Message must be at least 10 characters'),
})

type ContactFormData = z.infer<typeof contactSchema>

export default function ContactPage() {
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ContactFormData>({
    resolver: zodResolver(contactSchema),
  })

  const onSubmit = async (data: ContactFormData) => {
    setIsLoading(true)
    setError(null)
    setSuccess(false)

    try {
      const response = await apiClient.createContact({
        name: data.name,
        email: data.email,
        subject: data.subject || '',
        message: data.message,
      })

      if (response.success) {
        setSuccess(true)
        reset()
        // Reset success message after 5 seconds
        setTimeout(() => setSuccess(false), 5000)
      }
    } catch (err: any) {
      setError(err.message || 'Failed to send message. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12 px-4 sm:px-6 lg:px-8">
        <div className="container mx-auto max-w-2xl">
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
              <Mail className="h-8 w-8" />
            </div>
            <h1 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-3">
              Get in Touch
            </h1>
            <p className="text-lg text-secondary-600">
              Have a question or feedback? We'd love to hear from you. Send us a message and we'll respond as soon as possible.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle className="text-2xl">Contact Us</CardTitle>
              <CardDescription>Fill out the form below and we'll get back to you soon</CardDescription>
            </CardHeader>
            <form onSubmit={handleSubmit(onSubmit)}>
              <CardContent className="space-y-4">
                {error && (
                  <div className="p-3 rounded-lg bg-error-50 border border-error-200 text-error-600 text-sm">
                    {error}
                  </div>
                )}

                {success && (
                  <div className="p-3 rounded-lg bg-success-50 border border-success-500/20 text-success-600 text-sm flex items-center">
                    <CheckCircle2 className="h-4 w-4 mr-2" />
                    Thank you for your message! We'll get back to you soon.
                  </div>
                )}

                <Input
                  label="Full Name"
                  type="text"
                  placeholder="John Doe"
                  {...register('name')}
                  error={errors.name?.message}
                />

                <Input
                  label="Email"
                  type="email"
                  placeholder="you@example.com"
                  {...register('email')}
                  error={errors.email?.message}
                />

                <Input
                  label="Subject (Optional)"
                  type="text"
                  placeholder="What's this about?"
                  {...register('subject')}
                  error={errors.subject?.message}
                  helperText="Brief description of your inquiry"
                />

                <Textarea
                  label="Message"
                  placeholder="Tell us what's on your mind..."
                  rows={6}
                  {...register('message')}
                  error={errors.message?.message}
                  helperText="Please provide as much detail as possible"
                />
              </CardContent>
              <CardFooter className="flex flex-col space-y-4">
                <Button type="submit" className="w-full" isLoading={isLoading}>
                  <Send className="mr-2 h-4 w-4" />
                  Send Message
                </Button>
                <p className="text-sm text-center text-secondary-600">
                  We typically respond within 24-48 hours
                </p>
              </CardFooter>
            </form>
          </Card>

          <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card className="p-6">
              <h3 className="font-semibold text-secondary-900 mb-2">Email Us</h3>
              <p className="text-sm text-secondary-600">
                For general inquiries, reach out to us at{' '}
                <a href="mailto:support@joinus.io" className="text-primary-600 hover:text-primary-700">
                  support@joinus.io
                </a>
              </p>
            </Card>
            <Card className="p-6">
              <h3 className="font-semibold text-secondary-900 mb-2">Response Time</h3>
              <p className="text-sm text-secondary-600">
                We aim to respond to all messages within 24-48 hours during business days.
              </p>
            </Card>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  )
}

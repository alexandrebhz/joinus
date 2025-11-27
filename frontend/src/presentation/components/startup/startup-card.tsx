'use client'

import Link from 'next/link'
import Image from 'next/image'
import { Startup } from '@/domain/entities/startup.entity'
import { Card, CardContent, CardFooter } from '../ui/card'
import { Button } from '../ui/button'
import { Building2, MapPin, Users, Briefcase } from 'lucide-react'

interface StartupCardProps {
  startup: Startup
}

export function StartupCard({ startup }: StartupCardProps) {
  return (
    <Card hover className="h-full flex flex-col">
      <CardContent className="flex-1 pt-6">
        <div className="flex items-start space-x-4 mb-4">
          {startup.logoUrl ? (
            <div className="relative h-16 w-16 rounded-lg overflow-hidden flex-shrink-0 border border-secondary-200">
              <Image
                src={startup.logoUrl}
                alt={startup.name}
                fill
                className="object-cover"
              />
            </div>
          ) : (
            <div className="h-16 w-16 rounded-lg bg-primary-100 flex items-center justify-center flex-shrink-0">
              <Building2 className="h-8 w-8 text-primary-600" />
            </div>
          )}
          <div className="flex-1 min-w-0">
            <Link href={`/startups/${startup.slug}`}>
              <h3 className="text-lg font-semibold text-secondary-900 hover:text-primary-600 transition-colors mb-1">
                {startup.name}
              </h3>
            </Link>
            <p className="text-sm text-secondary-600 line-clamp-2">{startup.description}</p>
          </div>
        </div>

        <div className="space-y-2">
          <div className="flex items-center text-sm text-secondary-600">
            <MapPin className="h-4 w-4 mr-2 flex-shrink-0" />
            {startup.location}
          </div>
          <div className="flex items-center text-sm text-secondary-600">
            <Building2 className="h-4 w-4 mr-2 flex-shrink-0" />
            {startup.industry}
          </div>
          <div className="flex items-center gap-4 text-sm text-secondary-600">
            {startup.memberCount !== undefined && (
              <div className="flex items-center">
                <Users className="h-4 w-4 mr-1" />
                {startup.memberCount} members
              </div>
            )}
            {startup.jobCount !== undefined && startup.jobCount > 0 && (
              <div className="flex items-center">
                <Briefcase className="h-4 w-4 mr-1" />
                {startup.jobCount} jobs
              </div>
            )}
          </div>
        </div>
      </CardContent>

      <CardFooter>
        <Link href={`/startups/${startup.slug}`} className="w-full">
          <Button variant="outline" size="sm" className="w-full">
            View Profile
          </Button>
        </Link>
      </CardFooter>
    </Card>
  )
}


#!/usr/bin/env python

import webapp2
from google.appengine.api import app_identity
from google.appengine.api import mail
from conference import ConferenceApi

class SetAnnouncementHandler(webapp2.RequestHandler):
    def get(self):
        """Set Announcement in Memcache."""
        header = self.request.headers.get('X-AppEngine-Cron', None)
        if not header:
            raise ValueError('attempt to access cron handler directly, '
                             'missing custom App Engine header')
        """Set Announcement in Memcache."""
        # TODO 1
        # use _cacheAnnouncement() to set announcement in Memcache


app = webapp2.WSGIApplication([
    ('/crons/set_announcement', SetAnnouncementHandler),
], debug=True)
// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package admin

import (
        "html/template"
        "net/http"
        "time"

        "appengine"
        "appengine/datastore"
        "appengine/user"
)

// [START announcement_struct]
type Announcement struct {
        Author  string
        Content string
        Date    time.Time
}
// [END announcement_struct]

func init() {
        http.HandleFunc("/", root)
        http.HandleFunc("/announce", announce)
}

// announcementKey returns the key used for all announcement entries.
func announcementKey(c appengine.Context) *datastore.Key {
        // The string "default_feed" here could be varied to have multiple feeds.
        return datastore.NewKey(c, "Announcement", "default_feed", 0, nil)
}

// [START func_root]
func root(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        // Ancestor queries, as shown here, are strongly consistent with the High
        // Replication Datastore. Queries that span entity groups are eventually
        // consistent. If we omitted the .Ancestor from this query there would be
        // a slight chance that Greeting that had just been written would not
        // show up in a query.
        // [START query]
        q := datastore.NewQuery("Announcement").Ancestor(announcementKey(c)).Order("-Date").Limit(10)
        // [END query]
        // [START getall]
        announcements := make([]Announcement, 0, 10)
        if _, err := q.GetAll(c, &announcements); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        // [END getall]
        if err := feedTemplate.Execute(w, announcements); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}
// [END func_root]

var feedTemplate = template.Must(template.ParseFiles("index.html"))

// [START func_announce]
func announce(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        a := Announcement{
                Content: r.FormValue("content"),
                Date:    time.Now(),
        }
        if u := user.Current(c); u != nil {
                a.Author = u.String()
        }
        // We set the same parent key on every Announcement entity to ensure each 
        // Announcement is in the same entity group. Queries across the single entity 
        // group will be consistent. However, the write rate to a single entity group
        // should be limited to ~1/second.
        key := datastore.NewIncompleteKey(c, "Announcement", announcementKey(c))
        _, err := datastore.Put(c, key, &a)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}
// [END func_announce]

# NotiGoCMS

![CRON Update Job](https://github.com/lvntbkdmr/NotiGoCMS/workflows/CRON%20Update%20Job/badge.svg?branch=master)

This program is responsible for fetching your Notion pages and converting them into markdown files along with their static contents
Converted markdown files and their static contents are then pushed to your original Hugo blog git repo
Then, you only need to trigger "hugo" command to rebuild the required html files (If netlify is being used, this is done automatically)

kjk/notionapi is being used as api to fetch Notion data and kjk/blog is being referenced for conversion

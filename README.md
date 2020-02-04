# Forensic Files API

I'm building this project to learn/hone a whole bunch of skills. I'm planning on doing the following
things and I'm not sure if everything will shake out the way I want.

## Tasks

- [x] Extract Forensic Files episode links from YouTube (output goes into JSON file)
- [x] Loop through JSON file and download the videos using [youtube-dl](https://youtube-dl.org)
- [ ] Extract audio from video downloads
- [ ] Send audio to a speech-to-text service to transcribe episodes
- [ ] Use some kind of speech or AI parsing thing to categorize episode scripts
- [ ] Structure output of that in some way (maybe load into a graph database like [neo4j](http://neo4j.com/)?)
- [ ] Build interface to the database that allows you to search for stuff with keywords

## Goals

- Improve my Go skills
- Sharpen my JavaScript skills
- Learn how to use [ngrok](https://ngrok.com) with callback URLs (needed for speech-to-text service)
- Learn how to use text analysis tools
- Learn how to setup/interact with a graph database
- Learn how to setup a GraphQL backend

## Why Forensic Files?

I love Forensic Files and I have seen several episodes in which a detective or forensic specialist solves a case from watching another episode of Forensic Files (looking at you, diatoms)!

It would be neat to build a tool that allows you to search for stuff based on keywords (although it may be a little morbid).
[New episodes are coming out in February, 2020](https://nerdist.com/article/forensic-files-new-episodes-2020/) and if I got
props for helping to solve a murder on the show, I could die completely at peace (hopefully not from an [overdose of succinylcholine](https://www.imdb.com/title/tt1627318/?ref_=ttep_ep4)). 

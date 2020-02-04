/**
 * This is a simple script that I ran in Chrome DevTools to get the YouTube
 * video details I needed. I would do a search for "forensic files season X",
 * select a collection (thanks FilmRise), and run this script to log a JSON
 * string of all of the episodes in that season to the console. I just copied
 * and pasted that into `/assets/youtube-links.json`. Easy peasy!
 */

var endpointSelector = "a.yt-simple-endpoint.style-scope.ytd-playlist-panel-video-renderer";

var endpointAnchors = document.querySelectorAll(endpointSelector);

var videos = [];

endpointAnchors.forEach(endpointAnchor => {
    var videoTitleSpan = endpointAnchor.querySelector("#video-title");
    videos.push({ name: videoTitleSpan.title, url: endpointAnchor.href });
});

console.log(JSON.stringify(videos));

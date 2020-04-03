"use strict";

const path = require("path");
const fse = require("fs-extra");
const ytLinks = require("../assets/youtube-links.json");

cleanupEpisodes();

/**
 * One-off script to create a JSON file that I can parse to get the right
 * episode title when referring to file names and such.
 */
async function cleanupEpisodes() {
  const episodesBySeason = {};
  for (const [season, episodes] of Object.entries(ytLinks)) {
    episodesBySeason[season] = episodes.map(parseFields);
  }

  await fse.writeJSON(
    path.join("..", "assets", "episodes.json"),
    episodesBySeason,
    { spaces: 2 },
  );
}

function parseFields(episode) {
  const { name, url } = episode;
  const [seasonNumber, episodeNumber, title] = name.split(" | ");
  return {
    season: +seasonNumber.toString().replace("Season ", ""),
    episode: +episodeNumber.toString().replace("Episode ", ""),
    title: cleanupTitle(title),
    url,
  };
}

function cleanupTitle(title) {
  const chars = title.split("");
  const validChars = [];
  for (const char of chars) {
    if (/^[a-zA-Z]+$/.test(char)) {
      validChars.push(char);
    } else if (char === " ") {
      validChars.push("-");
    }
  }
  return validChars.join("").toLowerCase();
}

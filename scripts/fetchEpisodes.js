"use strict";

const path = require("path");
const qs = require("querystring");
const fetch = require("node-fetch");
const fse = require("fs-extra");

// Make sure you supply this in the `.env` file!
const API_KEY = process.env.OMDB_API_KEY;

fetchAllSeasons();

/**
 * Fetches all seasons of Forensic Files from the OMDb API and writes the
 * results to a JSON file in `/assets`.
 */
async function fetchAllSeasons() {
  const episodesBySeason = {};

  for (let season = 1; season < 15; season += 1) {
    console.log(`Fetching season ${season}...`);
    const seasonKey = season.toString().padStart(2, "0");
    episodesBySeason[seasonKey] = await fetchSeason(season);
  }

  const outputFilePath = path.resolve(__dirname, "..", "assets", "episodes.json");
  await fse.writeJSON(outputFilePath, episodesBySeason, { spaces: 2 });
  console.log("Done!");
}

async function fetchSeason(seasonNumber) {
  const queryString = qs.stringify({
    t: "Forensic Files",
    Season: seasonNumber,
    apikey: API_KEY
  });

  const response = await fetch(`http://www.omdbapi.com/?${queryString}`);
  const result = await response.json();
  return result.Episodes;
}

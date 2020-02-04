"use strict";

const path = require("path");
const fse = require("fs-extra");

const videosPath = path.resolve(__dirname, "..", "assets", "videos");

renameAllVideos();

/**
 * I decided to change the naming scheme for videos, so I whipped up this script
 * to knock that out really quick. I could have did it in Go, but I didn't want
 * to spend too much time on it.
 */
async function renameAllVideos() {
  for (let season = 1; season < 15; season += 1) {
    await renameVideosForSeason(season);
  }
}

async function renameVideosForSeason(season) {
  const seasonDirPath = path.join(videosPath, `season-${season}`);
  const videoFileNames = await fse.readdir(seasonDirPath);

  for (const videoFileName of videoFileNames) {
    const sourcePath = path.join(seasonDirPath, videoFileName);
    const seasonPrefix = season.toString().padStart(2, "0");
    const reCasedVideoFileName = videoFileName
      .toLowerCase()
      .replace(/ /gi, "-")
      .replace(/'/gi, "");
    const newName = seasonPrefix.concat("-", reCasedVideoFileName);
    const targetPath = path.join(seasonDirPath, newName);

    await fse.rename(sourcePath, targetPath);
  }
}

"use strict";

const path = require("path");
const fse = require("fs-extra");

const assetsPath = path.join(__dirname, "..", "assets");

createSeasonDirs("recognitions");

/**
 * Creates 14 directories in the specified parent directory with the "season-"
 * prefix.
 */
async function createSeasonDirs(parentName) {
  const parentDirPath = path.join(assetsPath, parentName);

  for (let i = 1; i <= 14; i++) {
    const seasonDirPath = path.join(parentDirPath, `season-${i}`);
    await fse.mkdirp(seasonDirPath);
  }
}

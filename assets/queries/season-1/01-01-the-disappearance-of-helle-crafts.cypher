CREATE (TheDisappearanceOfHelleCrafts:Episode {title:'The Disappearance of Helle Crafts', season: 1, episode: 1})
CREATE (Homicide:Event {name: 'Homicide'})
CREATE (WoodChipper:Object {name: 'Wood Chipper'})
CREATE (HelleCrafts:Person {name: 'Helle Crafts'})
CREATE (RichardCrafts:Person {name: 'Richard Crafts'})

CREATE
  (HelleCrafts)-[:APPEARED_IN]->(TheDisappearanceOfHelleCrafts),
  (RichardCrafts)-[:APPEARED_IN]->(TheDisappearanceOfHelleCrafts),
  (WoodChipper)-[:APPEARED_IN]->(TheDisappearanceOfHelleCrafts),
  (Homicide)-[:APPEARED_IN]->(TheDisappearanceOfHelleCrafts),
  (RichardCrafts)-[:USED]->(WoodChipper),
  (RichardCrafts)-[:PERPETRATED]->(Homicide),
  (HelleCrafts)-[:VICTIM_OF]->(Homicide)

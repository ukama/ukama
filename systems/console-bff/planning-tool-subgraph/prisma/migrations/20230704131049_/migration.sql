/*
  Warnings:

  - Made the column `draftId` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `north` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `west` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `east` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `south` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `url` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `populationUrl` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `populationCovered` on table `Site` required. This step will fail if there are existing NULL values in that column.
  - Made the column `totalBoxesCovered` on table `Site` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE "Site" ALTER COLUMN "draftId" SET NOT NULL,
ALTER COLUMN "north" SET NOT NULL,
ALTER COLUMN "west" SET NOT NULL,
ALTER COLUMN "east" SET NOT NULL,
ALTER COLUMN "south" SET NOT NULL,
ALTER COLUMN "url" SET NOT NULL,
ALTER COLUMN "populationUrl" SET NOT NULL,
ALTER COLUMN "populationCovered" SET NOT NULL,
ALTER COLUMN "totalBoxesCovered" SET NOT NULL;

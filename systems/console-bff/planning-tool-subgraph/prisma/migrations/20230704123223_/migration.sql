-- CreateTable
CREATE TABLE "Location" (
    "id" TEXT NOT NULL,
    "lat" TEXT NOT NULL,
    "lng" TEXT NOT NULL,
    "address" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "Link" (
    "id" TEXT NOT NULL,
    "siteA" TEXT NOT NULL,
    "siteB" TEXT NOT NULL,
    "draftId" TEXT
);

-- CreateTable
CREATE TABLE "Site" (
    "id" TEXT NOT NULL,
    "height" INTEGER NOT NULL,
    "solarUptime" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "apOption" TEXT NOT NULL,
    "draftId" TEXT,
    "isSetlite" BOOLEAN NOT NULL,
    "locationId" TEXT NOT NULL,
    "north" DOUBLE PRECISION NOT NULL,
    "west" DOUBLE PRECISION NOT NULL,
    "east" DOUBLE PRECISION NOT NULL,
    "south" DOUBLE PRECISION NOT NULL,
    "url" TEXT NOT NULL,
    "populationUrl" TEXT NOT NULL,
    "populationCovered" DOUBLE PRECISION NOT NULL,
    "totalBoxesCovered" DOUBLE PRECISION NOT NULL
);

-- CreateTable
CREATE TABLE "Event" (
    "id" TEXT NOT NULL,
    "operation" TEXT NOT NULL,
    "value" TEXT NOT NULL,
    "draftId" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateTable
CREATE TABLE "Draft" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "lastSaved" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "Location_id_key" ON "Location"("id");

-- CreateIndex
CREATE UNIQUE INDEX "Link_id_key" ON "Link"("id");

-- CreateIndex
CREATE UNIQUE INDEX "Site_id_key" ON "Site"("id");

-- CreateIndex
CREATE UNIQUE INDEX "Event_id_key" ON "Event"("id");

-- CreateIndex
CREATE UNIQUE INDEX "Draft_id_key" ON "Draft"("id");

-- CreateIndex
CREATE INDEX "Draft_userId_idx" ON "Draft" USING HASH ("userId");

-- AddForeignKey
ALTER TABLE "Link" ADD CONSTRAINT "Link_draftId_fkey" FOREIGN KEY ("draftId") REFERENCES "Draft"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Site" ADD CONSTRAINT "Site_locationId_fkey" FOREIGN KEY ("locationId") REFERENCES "Location"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Site" ADD CONSTRAINT "Site_draftId_fkey" FOREIGN KEY ("draftId") REFERENCES "Draft"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Event" ADD CONSTRAINT "Event_draftId_fkey" FOREIGN KEY ("draftId") REFERENCES "Draft"("id") ON DELETE CASCADE ON UPDATE CASCADE;

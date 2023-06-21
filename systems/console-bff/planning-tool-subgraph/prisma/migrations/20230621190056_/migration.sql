-- CreateTable
CREATE TABLE "Link" (
    "id" TEXT NOT NULL,
    "data" TEXT NOT NULL,
    "linkWith" TEXT NOT NULL,
    "siteId" TEXT
);

-- CreateIndex
CREATE UNIQUE INDEX "Link_id_key" ON "Link"("id");

-- AddForeignKey
ALTER TABLE "Link" ADD CONSTRAINT "Link_siteId_fkey" FOREIGN KEY ("siteId") REFERENCES "Site"("id") ON DELETE CASCADE ON UPDATE CASCADE;

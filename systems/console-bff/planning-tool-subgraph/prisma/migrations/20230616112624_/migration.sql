-- DropForeignKey
ALTER TABLE "Draft" DROP CONSTRAINT "Draft_siteId_fkey";

-- AddForeignKey
ALTER TABLE "Draft" ADD CONSTRAINT "Draft_siteId_fkey" FOREIGN KEY ("siteId") REFERENCES "Site"("id") ON DELETE CASCADE ON UPDATE CASCADE;

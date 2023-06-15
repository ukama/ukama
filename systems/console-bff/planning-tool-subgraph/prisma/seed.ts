import { Prisma, PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

const draftData: Prisma.DraftCreateInput[] = [
  {
    name: "Draft 1",
    lastSaved: 1686142941,
    userId: "16fe5842-05dd-11ee-8254-5eeae0fe08fe",
    createdAt: "2023-03-01T00:00:00.000Z",
    updatedAt: "2023-03-02T00:00:00.000Z",
    site: {
      create: {
        name: "Site 1",
        height: 100,
        apOption: "ONE_TO_ONE",
        solarUptime: 2,
        isSetlite: false,
        location: {
          create: {
            lat: "1001.123",
            lng: "12421.213",
            address: "Address 1",
          },
        },
      },
    },
    events: {
      create: [
        {
          operation: "name",
          value: "Site 2",
          createdAt: "2023-03-02T00:00:00.000Z",
        },
      ],
    },
  },
];

async function main() {
  console.log(`Start seeding ...`);
  for (const d of draftData) {
    const draft = await prisma.draft.create({
      data: d,
    });
    console.log(`Created data with id: ${draft.id}`);
  }
  console.log(`Seeding finished.`);
}

main()
  .then(async () => {
    await prisma.$disconnect();
  })
  .catch(async e => {
    console.error(e);
    await prisma.$disconnect();
    process.exit(1);
  });

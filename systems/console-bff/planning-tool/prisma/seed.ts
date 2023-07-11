import { Prisma, PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

const draftData: Prisma.DraftCreateInput[] = [
  {
    name: "Draft 1",
    userId: "a9a3dc45-fe06-43d6-b148-7508c9674627",
    createdAt: "2023-03-01T00:00:00.000Z",
    updatedAt: "2023-03-02T00:00:00.000Z",
    lastSaved: 1686142941,
    links: {
      create: [
        {
          id: "a50bb1a2-74a3-497c-bb98-b8fa9960c9a9",
          siteA: "a50bb7a3-74a3-497c-bb98-b8fa9960c9a9",
          siteB: "a50bb7a3-74a3-497c-bb99-b8fa8960c9a9",
        },
      ],
    },
    sites: {
      create: [
        {
          id: "a50bb7a3-74a3-497c-bb98-b8fa9960c9a9",
          name: "Site 1",
          height: 100,
          status: "up",
          apOption: "ONE_TO_ONE",
          solarUptime: 95,
          isSetlite: false,
          east: 0,
          west: 0,
          north: 0,
          south: 0,
          url: "",
          populationUrl: "",
          populationCovered: 0,
          totalBoxesCovered: 0,
          location: {
            create: {
              id: "a50bc7a3-74a3-497c-bb98-b8fa9960c9a9",
              lat: "-6.114665854",
              lng: "22.501076155",
              address: "Address 1",
            },
          },
        },
        {
          id: "a50bb7a3-74a3-497c-bb99-b8fa8960c9a9",
          name: "Site 2",
          height: 90,
          status: "up",
          apOption: "ONE_TO_TWO",
          solarUptime: 90,
          isSetlite: true,
          east: 0,
          west: 0,
          north: 0,
          south: 0,
          url: "",
          populationUrl: "",
          populationCovered: 0,
          totalBoxesCovered: 0,
          location: {
            create: {
              id: "a50bd7a3-74a3-497c-bb98-b8fa8960c9a9",
              lat: "-4.787720855",
              lng: "26.675040770",
              address: "Address 1",
            },
          },
        },
      ],
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

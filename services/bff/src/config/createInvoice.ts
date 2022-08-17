import fs from "fs";
import PDFDocument from "pdfkit";

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
const createInvoice = (subTotal: number, total: number, path: any) => {
    const doc: any = new PDFDocument({ size: "A4", margin: 50 });
    doc.font("Helvetica")
        .fontSize(16)
        .text(subTotal, 200, 50, { align: "right" })
        .moveDown(2)
        .text(total);
    doc.end();
    doc.pipe(fs.createWriteStream(path));
};

export { createInvoice };

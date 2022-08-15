import PDFDocument from "pdfkit";
const invoice = {
    shipping: {
        name: "John Doe",
        address: "1234 Main Street",
        city: "San Francisco",
        state: "CA",
        country: "US",
        postal_code: 94111,
    },
    items: [
        {
            item: "TC 100",
            description: "Toner Cartridge",
            quantity: 2,
            amount: 6000,
        },
        {
            item: "USB_EXT",
            description: "USB Cable Extender",
            quantity: 1,
            amount: 2000,
        },
    ],
    subtotal: 8000,
    // eslint-disable-next-line sort-keys
    paid: 0,
    invoice_nr: 1234,
};
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
const buildPDF = (dataCallback: any, endCallback: any) => {
    const doc = new PDFDocument({ bufferPages: true, font: "Courier" });
    doc.on("data", dataCallback);
    doc.on("end", endCallback);
    function generateHr(doc: any, y: any) {
        doc.strokeColor("grey")
            .lineWidth(1)
            .moveTo(50, y)
            .lineTo(550, y)
            .stroke();
    }
    function generateFooter(doc: any) {
        doc.fontSize(10).text(
            "Payment is due within 15 days. Thank you for your business.",
            50,
            780,
            { align: "center", width: 500 }
        );
    }

    function generateInvoiceTable(doc: any, invoice: any) {
        let i;
        const invoiceTableTop = 330;

        doc.font("Helvetica-Bold");
        generateTableRow(
            doc,
            invoiceTableTop,
            "Item",
            "Description",
            "Unit Cost",
            "Quantity",
            "Line Total"
        );
        generateHr(doc, invoiceTableTop + 20);
        doc.font("Helvetica");

        for (i = 0; i < invoice.items.length; i++) {
            const item = invoice.items[i];
            const position = invoiceTableTop + (i + 1) * 30;
            generateTableRow(
                doc,
                position,
                item.item,
                item.description,
                JSON.stringify(item.amount / item.quantity),
                item.quantity,
                JSON.stringify(item.amount)
            );

            generateHr(doc, position + 20);
        }

        const subtotalPosition = invoiceTableTop + (i + 1) * 30;
        generateTableRow(
            doc,
            subtotalPosition,
            "",
            "",
            "Subtotal",
            "",
            JSON.stringify(invoice.subtotal)
        );

        const paidToDatePosition = subtotalPosition + 20;
        generateTableRow(
            doc,
            paidToDatePosition,
            "",
            "",
            "Paid To Date",
            "",
            JSON.stringify(invoice.paid)
        );

        const duePosition = paidToDatePosition + 25;
        doc.font("Helvetica-Bold");
        generateTableRow(
            doc,
            duePosition,
            "",
            "",
            "Balance Due",
            "",
            JSON.stringify(invoice.subtotal - invoice.paid)
        );
        doc.font("Helvetica");
    }
    function generateTableRow(
        doc: any,
        y: any,
        item: any,
        description: any,
        unitCost: any,
        quantity: any,
        lineTotal: any
    ) {
        doc.fontSize(10)
            .text(item, 50, y)
            .text(description, 150, y)
            .text(unitCost, 280, y, { width: 90, align: "right" })
            .text(quantity, 370, y, { width: 90, align: "right" })
            .text(lineTotal, 0, y, { align: "right" });
    }

    // doc
    //   .image(
    //     'data:image/png;base64,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',
    //     25,
    //     25,
    //     { width: 150 }
    //   )
    //   .fillColor('#000')
    //   .fontSize(10)
    //   .font('Helvetica')
    //   .text(`Invoice ID    : 98879798`, { align: 'right' })
    //   .text(`Invoice number: 8782022`, { align: 'right' })
    //   .text(`Date          : €30, `, {
    //     align: 'right'
    //   })
    //   .text(`Network name  :     brackley, `, {
    //     align: 'right'
    //   })
    //   .text(`Status        :    Active , `, {
    //     align: 'right'
    //   })
    //   .text(`Country        :    DRC`, { align: 'right' })
    // doc.moveDown(2)
    // doc.fill('#021c27').text(`Bill from       Paid by  `, {
    //   align: 'right'
    // })

    // doc.moveDown()
    // const _kPAGE_BEGIN = 25
    // const _kPAGE_END = 580
    // //  [COMMENT] Draw a horizontal line.
    // doc.moveTo(_kPAGE_BEGIN, 200).lineTo(_kPAGE_END, 200).stroke()
    // doc.text(`Memo: jkjhkjhkjh`, 50, 210)
    // doc.moveTo(_kPAGE_BEGIN, 250).lineTo(_kPAGE_END, 250).stroke()

    // doc
    //   .fill('black')
    //   .font('Helvetica-Bold')
    //   .fontSize(14)
    //   .text('Ukama     CustomerName', {
    //     align: 'right'
    //   })
    // doc.moveDown(4)
    const customerInformationTop = 80;
    doc.image(
        "data:image/png;base64,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",
        25,
        25,
        { width: 100 }
    )
        .fontSize(10)
        .font("Helvetica")
        .fill("grey")
        .text("Invoice ID", 19, customerInformationTop, {
            width: 410,
            align: "right",
        })
        .font("Helvetica")
        .fill("black")
        .text(JSON.stringify(invoice.invoice_nr), 150, customerInformationTop, {
            width: 400,
            align: "right",
        })
        .fill("grey")
        .text("Invoice number", 55, customerInformationTop + 15, {
            width: 400,
            align: "right",
        })
        .fill("black")
        .text(
            JSON.stringify(invoice.subtotal - invoice.paid),
            150,
            customerInformationTop + 15,
            { width: 400, align: "right" }
        )
        .font("Helvetica")
        .fill("grey")
        .text("Date", 8, 110, { align: "right", width: 400 })
        .fill("black")
        .text("July 11 2022", 150, customerInformationTop + 30, {
            width: 400,
            align: "right",
        })
        .font("Helvetica")
        .fill("grey")
        .text("Network name", 50, 125, { align: "right", width: 400 })
        .fill("black")
        .text("Joe’s network", 150, 125, {
            width: 400,
            align: "right",
        })
        .font("Helvetica")
        .fill("grey")
        .text("Bill from", 25, 180, { align: "right", width: 400 })
        .font("Helvetica")
        .fill("grey")
        .text("Paid by", 125, 180, { align: "right", width: 400 })
        .font("Helvetica-Bold")
        .fill("black")
        .text("Ukama", 20, 200, { align: "right", width: 400 })
        .text("Customer name ", 150, 200, { align: "right", width: 400 });
    doc.moveDown(10);
    generateHr(doc, 250);
    doc.moveDown(3)
        .font("Helvetica-Bold")
        .fill("grey")
        .fontSize(14)
        .text("DESCRIPTION", 50, 300, { align: "left", width: 400 })
        .moveDown(4)
        .fontSize(14)
        .fill("black")
        .text("Standard roaming", 50, 340, { align: "left", width: 400 })
        .fontSize(10)
        .text("USD $20", 150, 350, { align: "right", width: 400 })
        .fontSize(10)
        .fill("grey")
        .text("5 GB; $5/GB; from June 11 2022 to July 11 2022 ", 52, 370, {
            align: "left",
            width: 400,
        })
        .fontSize(10)
        .fill("black")
        .text("Total", 50, 400, { align: "right", width: 400 })
        .text("$20", 150, 400, { align: "right", width: 400 })
        .strokeColor("grey")
        .lineWidth(1)
        .moveTo(425, 420)
        .lineTo(550, 420)
        .stroke()
        .moveDown(2)
        .fontSize(15)
        .text("Amount due", 110, 450, { align: "right", width: 400 })
        .text("$20", 150, 450, { align: "right", width: 400 });
    // doc
    //   .fill('grey')
    //   .text('Bill from        Paid by', {
    //     align: 'right'
    //   })
    //   .moveDown()
    // doc
    //   .fill('black')
    //   .text('Ukama        Casinga', 100, 170, {
    //     align: 'right'
    //   })
    //   .moveDown(10)

    // doc.fillColor('#444444').fontSize(20).text('Invoice', 50, 160)

    // generateHr(doc, 185)
    // doc
    //   .fontSize(10)
    //   .text('Invoice Number:', 50, customerInformationTop)
    //   .font('Helvetica-Bold')
    //   .text(JSON.stringify(invoice.invoice_nr), 150, customerInformationTop)
    //   .font('Helvetica')
    //   .text('Invoice Date:', 50, customerInformationTop + 15)
    //   .text(JSON.stringify(new Date()), 150, customerInformationTop + 15)
    //   .text('Balance Due:', 50, customerInformationTop + 30)
    //   .text(
    //     JSON.stringify(invoice.subtotal - invoice.paid),
    //     150,
    //     customerInformationTop + 30
    //   )

    //   .font('Helvetica-Bold')
    //   .text(invoice.shipping.name, 300, customerInformationTop)
    //   .font('Helvetica')
    //   .text(invoice.shipping.address, 300, customerInformationTop + 15)
    //   .text(
    //     invoice.shipping.city +
    //       ', ' +
    //       invoice.shipping.state +
    //       ', ' +
    //       invoice.shipping.country,
    //     300,
    //     customerInformationTop + 30
    //   )

    // generateHr(doc, 252)
    // generateInvoiceTable(doc, invoice)
    generateFooter(doc);

    doc.end();
};

export { buildPDF };

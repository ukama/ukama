import crypto from "crypto";

const generateTokenFromIccid = (iccid: string, key: string) => {
    const iccidEnvelope = {
        ICCID: iccid,
    };

    const tokenJson = JSON.stringify(iccidEnvelope);
    const encrypted = encrypt(tokenJson, key);

    return encrypted;
};

const encrypt = (t: string, key: string) => {
    if (key.length !== 32) {
        throw new Error("Key must be 32 bytes");
    }

    const cipher = crypto.createCipheriv(
        "aes-256-gcm",
        Buffer.from(key),
        Buffer.alloc(32)
    );

    const encrypted = Buffer.concat([cipher.update(t, "utf8"), cipher.final()]);
    const iv = cipher.getAuthTag();

    const result = Buffer.concat([iv, encrypted]);

    return result.toString("base64");
};

export default generateTokenFromIccid;

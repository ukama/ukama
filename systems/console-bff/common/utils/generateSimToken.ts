import crypto from "crypto";

const generateTokenFromIccid = (iccid: string, key: string) => {
  const iccidEnvelope = {
    ICCID: iccid,
  };

  const tokenJson = JSON.stringify(iccidEnvelope);
  const encrypted = encrypt(tokenJson, key);

  return encrypted;
};

const encrypt = (plaintext: string, key: string) => {
  if (key.length !== 32) {
    throw new Error("Key must be 32 bytes");
  }

  const iv = crypto.randomBytes(12);
  const cipher = crypto.createCipheriv("aes-256-gcm", Buffer.from(key), iv);

  const encrypted = Buffer.concat([
    cipher.update(plaintext, "utf8"),
    cipher.final(),
  ]);
  const tag = cipher.getAuthTag();

  const result = Buffer.concat([iv, encrypted, tag]);

  return result.toString("base64");
};

export default generateTokenFromIccid;

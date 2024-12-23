import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { createHash } from 'crypto';
import { Buffer } from 'buffer'; 
import { createGzip, Gzip } from 'zlib';
import { promisify } from 'util';
import { randomBytes, createCipheriv, createDecipheriv } from 'crypto';
import pako from 'pako';
import CryptoJS from 'crypto-js';

const secretKey = process.env.PASSWORD_SECRET

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
} 

export function formatNumber(num: number, decimals: number = 2): string {
    return num.toLocaleString('en-US', {
      minimumFractionDigits: decimals,
      maximumFractionDigits: decimals,
    });
  }
 
  
  // ฟังก์ชันสำหรับบีบอัดข้อมูล
  export const compressData = async (data: Buffer): Promise<Buffer> => {
      const gzip = createGzip();
      const buffer = await promisify(gzip.write.bind(gzip))(data);
      return buffer;
  };
  
  // ฟังก์ชันสำหรับเข้ารหัสข้อมูล
  export const encrypt = (data: Buffer, key: Buffer): string => {
      const iv = randomBytes(16); // สุ่ม IV
      const cipher = createCipheriv('aes-256-gcm', key, iv);
      const encrypted = Buffer.concat([cipher.update(data), cipher.final()]);
      return Buffer.concat([iv, encrypted]).toString('base64'); // รวม IV กับข้อมูลที่เข้ารหัส
  };
  
  // // ฟังก์ชันสำหรับถอดรหัสข้อมูลที่บีบอัด
  // export const decompressData = async (compressedData: Buffer): Promise<Buffer> => {
  //     const gunzip = promisify(require('zlib').gunzip);
  //     return await gunzip(compressedData);
  // };
  
  export function decompressGzip(compressedData: Uint8Array): Uint8Array | null {
    try {
        const decompressed = pako.inflate(compressedData);
        return decompressed;
    } catch (err) {
        console.error('Decompression failed:', err);
        return null;
    }
}

  // ฟังก์ชันสำหรับถอดรหัสข้อมูลที่เข้ารหัส
  export const decrypt = (encryptedData: string): Buffer => {
      const key = Buffer.from(secretKey || '', 'utf-8'); 
      const data = Buffer.from(encryptedData, 'base64');
      const iv = data.slice(0, 16); // แยก IV
      const encryptedText = data.slice(16); // ข้อมูลที่เข้ารหัส
      const decipher = createDecipheriv('aes-256-gcm', key, iv);
      return Buffer.concat([decipher.update(encryptedText), decipher.final()]);
  };

  export const Decrypt = (encryptedData:string):string => {
    const key = CryptoJS.enc.Utf8.parse(secretKey); // คีย์ 32 ไบต์

    try {
        // ถอดรหัส AES-GCM
        const decryptedData = CryptoJS.AES.decrypt(encryptedData, key, {
            mode: CryptoJS.mode.GCM,
            padding: CryptoJS.pad.NoPadding,
        });

        // แปลงข้อมูลกลับเป็นข้อความ
        return decryptedData.toString(CryptoJS.enc.Utf8);
      }   catch (error) {
        return 'Error processing received data:'+error;
    }
  }
// export function deCompressed(txt: string) : string {
	
//   const key = Buffer.from(secretKey || '', 'utf-8'); 
//   const hashedKey = createHash('sha256').update(key).digest();

//   // decryptedData,err := jwtn.Decrypt(msg.Payload,hashedKey[:])
// 	// 	if err != nil {
// 	// 		log.Fatalf("Could not decrypted: %v", err)
// 	// 	}
// 	// 	decompressedData,err := jwtn.DecompressData(decryptedData)
// 	// 	if err != nil {
// 	// 		log.Fatalf("Could not decompress: %v", err)
// 	// 	}

//   return ""
// }
 

import { kvsLocalStorage } from '@kvs/node-localstorage'

const storeInBucket = async (key: string, value: any) => {
 const storageKeyValue = await kvsLocalStorage({
  name: process.env.STORAGE_KEY,
  version: 1,
 })
 await storageKeyValue.set(key, value)
}

const retriveFromBucket = async (key: string): Promise<any> => {
 const storageKeyValue = await kvsLocalStorage({
  name: process.env.STORAGE_KEY,
  version: 1,
 })
 return await storageKeyValue.get(key)
}

const removeKeyFromBucket = async (key: string) => {
 const storageKeyValue = await kvsLocalStorage({
  name: process.env.STORAGE_KEY,
  version: 1,
 })
 await storageKeyValue.delete(key)
}

export { removeKeyFromBucket, retriveFromBucket, storeInBucket }

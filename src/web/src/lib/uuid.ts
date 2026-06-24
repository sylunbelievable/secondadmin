type UUIDCrypto = {
  randomUUID?: () => string
  getRandomValues?: Crypto['getRandomValues']
}

type DeviceStorage = Pick<Storage, 'getItem' | 'setItem'>

export function createUUID(cryptoSource: UUIDCrypto | undefined = globalThis.crypto) {
  if (typeof cryptoSource?.randomUUID === 'function') return cryptoSource.randomUUID()

  const bytes = new Uint8Array(16)
  if (typeof cryptoSource?.getRandomValues === 'function') {
    cryptoSource.getRandomValues(bytes)
  } else {
    for (let index = 0; index < bytes.length; index++) bytes[index] = Math.floor(Math.random() * 256)
  }

  bytes[6] = (bytes[6] & 0x0f) | 0x40
  bytes[8] = (bytes[8] & 0x3f) | 0x80

  return [...bytes].map((byte, index) => {
    const hex = byte.toString(16).padStart(2, '0')
    return [4, 6, 8, 10].includes(index) ? `-${hex}` : hex
  }).join('')
}

export function getOrCreateDeviceId(storage: DeviceStorage = localStorage) {
  const existing = storage.getItem('deviceId')
  if (existing) return existing

  const id = createUUID()
  storage.setItem('deviceId', id)
  return id
}

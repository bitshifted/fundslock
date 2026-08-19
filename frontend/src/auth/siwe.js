 
import {createSIWEConfig, formatMessage,} from '@reown/appkit-siwe'
import { getNonce, verifySignature, getUserSession, refreshToken } from "./authenticate.js";
import { useAuthStore } from '@/stores/auth.js';


async function getSession(){
  const authStore = useAuthStore()
  try {
    const curToken = authStore.token
    try {
      let session = await getUserSession(curToken)
      console.log("session: " + session)
      if(!session) {
        console.log("Attempting to refresh token")
        const newToken = await refreshToken()
        console.log("New refresh token: " + newToken)
        authStore.setToken(newToken)
        if(newToken) {
          session = await getUserSession(newToken)
        }
      }
      return session
    } catch(error) {
      console.log(error)
      return null
    }
    
  } catch (error) {
    return null
  
  }
}

async function verifyMessage({ message, signature }){
  try {
    const isValid = await verifySignature(message, signature)

    return isValid
  } catch (error) {
    return false
  }
}

const siweConfig = createSIWEConfig({
    getMessageParams: async () => ({
    domain: window.location.host,
    uri: window.location.origin,
    chains: [1],
    statement: 'Please sign with your account',
  }),
  createMessage: ({ address, ...args }) => formatMessage(args, address),
  getNonce: async () => { 
    const nonce = getNonce()
    if (!nonce) {
      throw new Error('Failed to get nonce!')
    }
    return nonce
  },
  getSession,
  verifyMessage,
  signOut: async () => { //Example
    // Implement your Sign out function
  }
})

export { siweConfig }
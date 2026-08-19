import { BACKEND_URL } from '@/config/common.js'

const NONCE_URL = `${BACKEND_URL}/api/v1/auth/nonce`
const REFRESH_URL = `${BACKEND_URL}/api/v1/auth/refresh`
const VERIFY_URL = `${BACKEND_URL}/api/v1/auth/verify`
const SESSION_URL = `${BACKEND_URL}/api/v1/users/session`




async function getNonce() {
    const nonceRes = await fetch(NONCE_URL)
    const { nonce } = await nonceRes.json()
    console.log("nonce: " + nonce)
    return nonce
}

async function verifySignature(siweMessage, signature) {
    try {
        const response = await fetch(VERIFY_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message: siweMessage, signature:  signature }),
            credentials: 'include',
        })
        if (!response.ok) {
            throw new Error('Error verifying signature: ${response.status} ')
        }
        return true
    } catch (error) {
        console.log(error)
        return false
    }
}

async function refreshToken() {
    const response = await fetch(REFRESH_URL, {
        credentials: 'include',
        method: 'GET',
    })
    if(response.status === 401) {
        console.log("Failed to refresh token: ${response.status}")
        return null
    }
    if (!response.ok) {
        throw new Error('Error refresing token: ${response.status} ')
    }
    const { access_token } = await response.json()
    return access_token
   
}

async function getUserSession(jwtToken) {
    console.log("Fetching user session")
    if(!jwtToken) {
        console.log("Token not present")
        return null
    }
    const response = await fetch(SESSION_URL, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${jwtToken}`,
        },
    })
    if(response.status === 401) {
        return null
    }
    if (!response.ok) {
        throw new Error('Error getting session: ${response.status} ')
    }
    const session = await response.json()
    return session
}

export { getNonce, refreshToken, verifySignature, getUserSession }

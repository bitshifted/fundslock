
import { defineStore } from "pinia"
import {ref} from "vue"


export const useAuthStore =  defineStore('auth', () => {
    const token = ref(null)

    function setToken(newValue) {
        token.value = newValue
    }

    function clearToken() {
        token.value = null
    }


    return { 
        token, setToken , clearToken
    }
    
})
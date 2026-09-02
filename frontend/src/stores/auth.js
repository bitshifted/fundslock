
import { defineStore } from "pinia"
import {computed, ref} from "vue"


export const useAuthStore =  defineStore('auth', () => {
    const token = ref(null)
    const isLoggedIn = computed(() => token != null)

    function setToken(newValue) {
        token.value = newValue
    }

    function clearToken() {
        token.value = null
    }


    return { 
        token, setToken , clearToken, isLoggedIn
    }
    
})
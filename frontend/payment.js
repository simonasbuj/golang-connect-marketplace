function paymentApp() {
    return {
        isSuccess: false,
        isError: false,
        appUrl: "exp://192.168.1.137:8081/--/(tabs)/profile",

        async init () {
            this.checkUrlParams()
        },

        checkUrlParams() {
            const params = new URLSearchParams(window.location.search)

            this.isSuccess = params.get("success") === "true"
            this.isError = params.get("error") === "true"
        }
    }
}
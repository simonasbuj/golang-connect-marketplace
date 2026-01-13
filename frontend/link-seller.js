function linkSellerApp() {
    return {
        isSuccess: false,
        isError: false,
        appUrl: "goconnectmarketplace://",

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
function oauthApp() {
    return {
        code: null,
        isError: false,
        isSuccess: false,
        copiedMessage: "",
        appUrl: "exp://192.168.1.137:8081/--/(auth)/oauth",

        async init () {
            this.getCodeParam()
        },

        getCodeParam() {
            const params = new URLSearchParams(window.location.search)

            code = params.get("code") || ''
            
            this.code = code

            if (code) {
                this.isSuccess = true
            } else {
                this.isError = true
            }
        },

        copyCode() {
            if (!this.code) return

            navigator.clipboard.writeText(this.code)
                .then(() => {
                    this.copiedMessage = 'Copied!'
                    setTimeout(() => this.copiedMessage = '', 1500)
                })
                .catch(() => {
                    this.copiedMessage = 'Failed to copy'
                    setTimeout(() => this.copiedMessage = '', 1500)
                })
        }
    }
}
function oauthApp() {
    return {
        code: "",
        isError: false,
        isSuccess: false,
        copiedMessage: "",
        appUrl: "goconnectmarketplace://(auth)/oauth-exchange?code={CODE}",

        async init () {
            this.getCodeParam()
        },

        getCodeParam() {
            const params = new URLSearchParams(window.location.search)

            code = params.get("code") || ''
            
            this.code = code

            if (code) {
                this.isSuccess = true
                this.appUrl = this.appUrl.replace("{CODE}", code)
            } else {
                this.isError = true
            }
        },

        copyCode() {
            if (!this.code) return

            if (navigator.clipboard && window.isSecureContext) {
                navigator.clipboard.writeText(this.code)
                .then(() => {
                    this.showCopied()
                })
                .catch(() => {
                    this.fallbackCopy()
                })
            } else {
                this.fallbackCopy()
            }
        },

        fallbackCopy() {
            const textarea = document.createElement('textarea')
            textarea.value = this.code
            textarea.setAttribute('readonly', '')
            textarea.style.position = 'absolute'
            textarea.style.left = '-9999px'

            document.body.appendChild(textarea)
            textarea.select()
            textarea.setSelectionRange(0, textarea.value.length)

            try {
                document.execCommand('copy')
                this.showCopied()
            } catch (e) {
                this.showFailed()
            }

            document.body.removeChild(textarea)
        },

        showCopied() {
            this.copiedMessage = 'Copied!'
            setTimeout(() => this.copiedMessage = '', 1500)
        },

        showFailed() {
            this.copiedMessage = 'Tap & hold to copy'
            setTimeout(() => this.copiedMessage = '', 2000)
            }
        }
}
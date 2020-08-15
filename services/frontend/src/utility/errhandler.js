const errhandler = {
    handle: function (error) {
        if (error.response) {
            return this._outer_error_handler(error.response)
        }
        return this._inner_error_handler(error)
    },

    _outer_error_handler(error) {
        let msg = ""
        const code = error.status

        switch (code) {
            case 401:
                msg = 'unauthorized'
                break
            case 403:
                msg = "unauthenticated"
                break
            case 500:
                msg = "internal server error"
                break
            case 503:
                msg = "service unavailable"
                break
            default:
                msg = error.data
        }
        return { message: msg, code: code }
    },

    _inner_error_handler(error) {
        return { message: error.message, code: null }
    }
}

export default errhandler
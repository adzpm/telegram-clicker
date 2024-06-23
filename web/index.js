let Application = Vue.createApp({
    data: function () {
        return {
            telegram_data: {},

            telegram_id: 777,
            tmp_btn_count: [1],
            user: {
                id: 0,
                telegram_id: 0,
                coins: 0,
                last_seen: 0,
            },
        }
    },

    methods: {
        Login() {
            let url = this.CurrentURL + '/login' + '?telegram_id=' + this.telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        Click() {
            let url = this.CurrentURL + '/click' + '?telegram_id=' + this.telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user = response.data
            }).catch(error => {
                console.log(error)
            })
        }
    },

    computed: {
        CurrentURL() {
            let url = window.location.href
            return url.substring(0, url.lastIndexOf('/'))
        },

        TelegramID() {
            if (window.Telegram.WebApp.initialDataUnsafe &&
                window.Telegram.WebApp.initialDataUnsafe.user &&
                window.Telegram.WebApp.initialDataUnsafe.user.id) {
                return window.Telegram.WebApp.initialDataUnsafe.user.id
            }

            return "123456789"
        },

        TelegramName() {
            let display_name = ""

            if (window.Telegram.WebApp.initialDataUnsafe &&
                window.Telegram.WebApp.initialDataUnsafe.user &&
                window.Telegram.WebApp.initialDataUnsafe.user.first_name) {
                display_name = window.Telegram.WebApp.initialDataUnsafe.user.first_name
            }

            if (window.Telegram.WebApp.initialDataUnsafe &&
                window.Telegram.WebApp.initialDataUnsafe.user &&
                window.Telegram.WebApp.initialDataUnsafe.user.last_name) {
                display_name += " " + window.Telegram.WebApp.initialDataUnsafe.user.last_name
            }

            return display_name ? display_name : "John Doe"
        }
    },

    mounted: function () {
        this.Login()
    },
})

const vm = Application.mount('#app')
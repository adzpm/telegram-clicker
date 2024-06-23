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
            return "123456789"
        },

        TelegramName() {
            return window.Telegram?.WebApp?.initialDataUnsafe?.user?.first_name
        }
    },

    mounted: function () {
        // copy telegram data to local variable
        this.Login()
    },
})

const vm = Application.mount('#app')
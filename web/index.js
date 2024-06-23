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
        Login(telegram_id) {
            let url = this.CurrentURL + '/login' + '?telegram_id=' + telegram_id

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
            return window.Telegram?.WebApp?.initDataUnsafe?.user?.id
        },

        TelegramName() {
            return window.Telegram?.WebApp?.initDataUnsafe?.user?.first_name
        }
    },

    mounted: function () {
        // copy telegram data to local variable
        while (window.Telegram.WebApp.initDataUnsafe === undefined) {
            setTimeout(() => {}, 10)
        }

        this.telegram_data = {...window.Telegram?.WebApp?.initDataUnsafe}
        console.log(this.telegram_data)

        this.Login(this.TelegramID)

    },
})

const vm = Application.mount('#app')
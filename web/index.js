let Application = Vue.createApp({
    data: function () {
        return {
            telegram_data: {},

            tmp_btn_count: [1],
            user_data: {
                coins: 0,
            },
        }
    },

    methods: {
        Login(telegram_id) {
            let url = this.CurrentURL + '/login' + '?telegram_id=' + telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user_data = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        Click(telegram_id) {
            let url = this.CurrentURL + '/click' + '?telegram_id=' + telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user_data = response.data
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
            return this.telegram_data?.user?.id ?? 9876543210
        },

        TelegramName() {
            let name = this.telegram_data?.user?.first_name

            if (this.telegram_data?.user?.last_name) {
                name += ' ' + this.telegram_data?.user?.last_name
            }

            return name ?? 'John Doe'
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
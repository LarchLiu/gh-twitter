<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>
    <div class="header">
        <a-button type="primary" @click="onExit">
            退出
        </a-button>
    </div>
    <a-calendar v-model="value">
        <template v-slot:dateCellRender="{ current: value }">
            <ul class="events">
                <li v-for="item in usersData"
                    :key="item.Profile.Name">
                    <a-badge status="success"
                    :text="item.Profile.Username" />
                </li>
            </ul>
        </template>
        <template v-slot:monthCellRender="{ current: value }">
            <div v-if="getMonthData(value)" class="notes-month">
                <section>{{ getMonthData(value) }}</section>
                <span>Backlog number</span>
            </div>
        </template>
    </a-calendar>
  </div>
</template>
<script>
import { ref,getCurrentInstance, reactive, onMounted, watch } from 'vue';
import userApi from '/@/api/twitter/user';


export default {
    setup () {

        const { ctx } = getCurrentInstance();

        let value =  '';
        const usersList = ref([])
        const usersData = ref([])

        const getListData = (value)=> {
            let listData
            switch (value.date()) {
                case 8:
                    listData = [
                        { type: 'warning', content: 'This is warning event.' },
                        { type: 'success', content: 'This is usual event.' },
                    ]
                    break
                case 10:
                    listData = [
                        { type: 'warning', content: 'This is warning event.' },
                        { type: 'success', content: 'This is usual event.' },
                        { type: 'error', content: 'This is error event.' },
                    ]
                    break
                case 15:
                    listData = [
                        { type: 'warning', content: 'This is warning event' },
                        {
                            type: 'success',
                            content: 'This is very long usual event。。....',
                        },
                        { type: 'error', content: 'This is error event 1.' },
                        { type: 'error', content: 'This is error event 2.' },
                        { type: 'error', content: 'This is error event 3.' },
                        { type: 'error', content: 'This is error event 4.' },
                    ]
                    break
                default:
            }
            return listData || []
        }

        const getMonthData = (value) =>{
            if (value.month() === 8) {
                return 1394
            }
        }

        const onExit = ()=>{
            ctx.$router.push({
                path: '/login'
            })
        }

        const getUserList = ()=> {
            userApi.getUsersData().then(res => {
                console.log(res)
                usersList.value = res.replace(/\s*/g,"").split(',')
            }).catch(e => {
                console.log(e)
                usersList.value = []
            })
        }

        const getUserTweets = (user) => {
            userApi.getTweetsData(user).then(data => {
                usersData.value.push(data)
                console.log(data)
            }).catch(err => {
                console.log(err)
            })
        }

        onMounted(() => {
          getUserList()
        })
        watch(usersList, () => {
            usersData.value = []
            for (let i = 0; i < usersList.value.length; i++) {
                getUserTweets(usersList.value[i])
            }
        })

        return {
            value,
            usersList,
            usersData,
            getListData,
            getMonthData,
            getUserList,
            onExit
        }
    }
}
</script>
<style scoped lang="less">
.events {
    list-style: none;
    margin: 0;
    padding: 0;
}
.events .ant-badge-status {
    overflow: hidden;
    white-space: nowrap;
    width: 100%;
    text-overflow: ellipsis;
    font-size: 12px;
}
.notes-month {
    text-align: center;
    font-size: 28px;
}
.notes-month section {
    font-size: 28px;
}
.header{
    text-align: right;
    padding: 10px;
}
</style>

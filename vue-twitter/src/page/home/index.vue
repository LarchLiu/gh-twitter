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
                <li v-for="item in getListData(value)"
                    :key="item.content">
                    <a-badge :status="item.type"
                    :text="item.content" />
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
import { ref,getCurrentInstance, reactive } from 'vue';
import jsonp from '/@/api/jsonp';
import { commonParams, options } from '/@/api/config';

import userApi from '/@/api/twitter/user';


export default {
    setup () {

        const { ctx } = getCurrentInstance();

        let value =  '';

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
        return {
            value,
            getListData,
            getMonthData,
            onExit
        }
    },
    created () {
      // userApi.getUsersData().then(res => {
      //   console.log(res)
      // }).catch(e => {
      //   console.log(e)
      // })
      // const url ="https://c.y.qq.com/v8/fcg-bin/fcg_myqq_toplist.fcg"  //跨域的url
      const url = "https://raw.githubusercontent.com/LarchLiu/twitter-json/main/raw/twitter/userList.txt"
      const data = Object.assign({},commonParams,{  //Object.assign()是ES6api是将后面的对象参数，写入{}中，去重，且如存在相同key后者覆盖前者。
        platform: 'h5',
        format: 'text',
        needNewCode: 1
      })
      jsonp(url, null, options).then(res => {
        console.log(res)
      }).catch(err => {
        console.log(err)
      })
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

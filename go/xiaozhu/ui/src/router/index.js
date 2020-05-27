import Vue from 'vue'
import Router from 'vue-router'
//
import tipswitch from '@/pages/tipswitch'
import friendswitch from '@/pages/friendswitch'
import water8switch from '@/pages/water8switch'

//

Vue.use(Router)

const router = new Router({
  base: '/',
  mode: 'hash',
  routes: [
    {
      path: "/",
      redirect: "/tipswitch",
    },
    {
      path: '/tipswitch',
      meta: { name: "打卡提醒管理" },
      component: tipswitch
    },
    {
      path: '/friendswitch',
      meta: { name: "朋友管理" },
      component: friendswitch
    },
    {
      path: '/water8switch',
      meta: { name: "8杯水提醒" },
      component: water8switch
    },



  ]
});

//
export default router;

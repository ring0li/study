<style lang="less" scoped>
.value-class {
  flex: none !important;
}
</style>
<template>
  <div>
    <van-radio-group :value="waterval" @change="onChange">
      <van-cell-group v-for="(item,idx) in flist" :key="idx">
        <van-cell
          :title="item.name"
          value-class="value-class"
          clickable
          :data-name="item.val"
          @click="onClick"
        >
          <van-radio :name="item.val" />
        </van-cell>
      </van-cell-group>
    </van-radio-group>
  </div>
</template>

<script>
import api from "./data";

export default {
  name: "water8",
  components: {},
  data() {
    return {
      waterval: 8,
      openid: this.$route.query.openid,
      flist: [
        { name: "8杯水", val: 8 },
        { name: "6杯水", val: 6 },
        { name: "4杯水", val: 4 },
        { name: "不提醒", val: 0 }
      ]
    };
  },
  mounted() {
    document.title = "8杯水提醒";

    this._getWaterTipData();
  },
  methods: {
    _getWaterTipData() {
      let _this = this;
      api.getWaterTipData({ openId: this.openid }).then(res => {
        if (res.errorCode == 0) {
          _this.waterval = res.data.Waterval;
        }
      });
    },
    _setWaterTipData() {
      let _this = this;
      api
        .setWaterTipData({ openId: this.openid, waterval: _this.waterval })
        .then(res => {
          if (res.errorCode == 0) {
            _this.$toast.success("设置成功");
          } else {
            _this.$toast(res.message);
          }
        });
    },
    onChange(event) {
      // this.radio = event.detail;
    },
    onClick(event) {
      this.waterval = parseInt(event.currentTarget.dataset.name);

      this._setWaterTipData();
    }
  }
};
</script>


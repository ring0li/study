<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>内部使用切记外传，后果自负</title>
    <!-- 引入 echarts.js -->
    <script src="/static/js/echarts.min.js"></script>
</head>
<body>



 <div id="div1" style="width: 100%;height:400px;"></div>
        <script type="text/javascript">
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('div1'));

            // 指定图表的配置项和数据
            var option = {
                title: {
                    text: ''
                },
                tooltip: {},
                legend: {
                    data:['']
                },
                xAxis: {
                    data: {{.x1}}
                },
                yAxis: {},
                series: [{
                    name: '',
                    type: 'bar',
                    data: {{.y1}}
                }]
            };

            // 使用刚指定的配置项和数据显示图表。
            myChart.setOption(option);
        </script>


    <!-- 为ECharts准备一个具备大小（宽高）的Dom -->
    <div id="div2" style="width: 100%;height:400px;"></div>
    <script type="text/javascript">
        // 基于准备好的dom，初始化echarts实例
        var myChart = echarts.init(document.getElementById('div2'));

        // 指定图表的配置项和数据
        var option = {
            title: {
                text: ''
            },
            tooltip: {},
            legend: {
                data:['']
            },
            xAxis: {
                data: {{.x2}}
            },
            yAxis: {},
            series: [{
                name: '',
                type: 'bar',
                data: {{.y2}}
            }]
        };

        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    </script>



 <!-- 为ECharts准备一个具备大小（宽高）的Dom -->
    <div id="div3" style="width: 100%;height:400px;"></div>
    <script type="text/javascript">
        // 基于准备好的dom，初始化echarts实例
        var myChart = echarts.init(document.getElementById('div3'));

        // 指定图表的配置项和数据
        var option = {
            title: {
                text: ''
            },
            tooltip: {},
            legend: {
                data:['']
            },
            xAxis: {
                data: {{.x3}}
            },
            yAxis: {},
            series: [{
                name: '',
                type: 'bar',
                data: {{.y3}}
            }]
        };

        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    </script>



    <div id="div4" style="width: 100%;height:400px;"></div>
    <script type="text/javascript">
        // 基于准备好的dom，初始化echarts实例
        var myChart = echarts.init(document.getElementById('div4'));

        // 指定图表的配置项和数据
        var option = {
            title: {
                text: ''
            },
            tooltip: {},
            legend: {
                data:['']
            },
            xAxis: {
                data: {{.x4}}
            },
            yAxis: {},
            series: [{
                name: '',
                type: 'bar',
                data: {{.y4}}
            }]
        };

        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    </script>



     <div id="div_guanzhu" style="width: 100%;height:400px;"></div>
        <script type="text/javascript">
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('div_guanzhu'));

            // 指定图表的配置项和数据
            var option = {
                title: {
                    text: '雪球关注数'
                },
                tooltip: {},
                legend: {
                    data:['']
                },
                xAxis: {
                    data: {{.guanzhuX}}
                },
                yAxis: {},
                series: [{
                    name: '',
                    type: 'bar',
                    data: {{.guanzhuY}}
                }]
            };

            // 使用刚指定的配置项和数据显示图表。
            myChart.setOption(option);
        </script>


</body>
</html>
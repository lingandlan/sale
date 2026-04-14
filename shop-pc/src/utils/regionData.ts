// 中国省市区三级联动数据
// 用于充值中心新建/编辑时的地区选择

export interface RegionOption {
  label: string
  value: string
  children?: RegionOption[]
}

export const regionData: RegionOption[] = [
  {
    label: '北京市', value: '北京市', children: [
      { label: '北京市', value: '北京市', children: [
        { label: '东城区', value: '东城区' },
        { label: '西城区', value: '西城区' },
        { label: '朝阳区', value: '朝阳区' },
        { label: '丰台区', value: '丰台区' },
        { label: '石景山区', value: '石景山区' },
        { label: '海淀区', value: '海淀区' },
        { label: '门头沟区', value: '门头沟区' },
        { label: '房山区', value: '房山区' },
        { label: '通州区', value: '通州区' },
        { label: '顺义区', value: '顺义区' },
        { label: '昌平区', value: '昌平区' },
        { label: '大兴区', value: '大兴区' },
        { label: '怀柔区', value: '怀柔区' },
        { label: '平谷区', value: '平谷区' },
        { label: '密云区', value: '密云区' },
        { label: '延庆区', value: '延庆区' },
      ] }
    ]
  },
  {
    label: '上海市', value: '上海市', children: [
      { label: '上海市', value: '上海市', children: [
        { label: '黄浦区', value: '黄浦区' },
        { label: '徐汇区', value: '徐汇区' },
        { label: '长宁区', value: '长宁区' },
        { label: '静安区', value: '静安区' },
        { label: '普陀区', value: '普陀区' },
        { label: '虹口区', value: '虹口区' },
        { label: '杨浦区', value: '杨浦区' },
        { label: '闵行区', value: '闵行区' },
        { label: '宝山区', value: '宝山区' },
        { label: '嘉定区', value: '嘉定区' },
        { label: '浦东新区', value: '浦东新区' },
        { label: '金山区', value: '金山区' },
        { label: '松江区', value: '松江区' },
        { label: '青浦区', value: '青浦区' },
        { label: '奉贤区', value: '奉贤区' },
        { label: '崇明区', value: '崇明区' },
      ] }
    ]
  },
  {
    label: '天津市', value: '天津市', children: [
      { label: '天津市', value: '天津市', children: [
        { label: '和平区', value: '和平区' },
        { label: '河东区', value: '河东区' },
        { label: '河西区', value: '河西区' },
        { label: '南开区', value: '南开区' },
        { label: '河北区', value: '河北区' },
        { label: '红桥区', value: '红桥区' },
        { label: '东丽区', value: '东丽区' },
        { label: '西青区', value: '西青区' },
        { label: '津南区', value: '津南区' },
        { label: '北辰区', value: '北辰区' },
        { label: '武清区', value: '武清区' },
        { label: '宝坻区', value: '宝坻区' },
        { label: '滨海新区', value: '滨海新区' },
      ] }
    ]
  },
  {
    label: '重庆市', value: '重庆市', children: [
      { label: '重庆市', value: '重庆市', children: [
        { label: '万州区', value: '万州区' },
        { label: '涪陵区', value: '涪陵区' },
        { label: '渝中区', value: '渝中区' },
        { label: '大渡口区', value: '大渡口区' },
        { label: '江北区', value: '江北区' },
        { label: '沙坪坝区', value: '沙坪坝区' },
        { label: '九龙坡区', value: '九龙坡区' },
        { label: '南岸区', value: '南岸区' },
        { label: '北碚区', value: '北碚区' },
        { label: '渝北区', value: '渝北区' },
        { label: '巴南区', value: '巴南区' },
        { label: '江津区', value: '江津区' },
        { label: '合川区', value: '合川区' },
        { label: '永川区', value: '永川区' },
      ] }
    ]
  },
  {
    label: '广东省', value: '广东省', children: [
      { label: '广州市', value: '广州市', children: [
        { label: '荔湾区', value: '荔湾区' },
        { label: '越秀区', value: '越秀区' },
        { label: '海珠区', value: '海珠区' },
        { label: '天河区', value: '天河区' },
        { label: '白云区', value: '白云区' },
        { label: '黄埔区', value: '黄埔区' },
        { label: '番禺区', value: '番禺区' },
        { label: '花都区', value: '花都区' },
        { label: '南沙区', value: '南沙区' },
        { label: '增城区', value: '增城区' },
      ] },
      { label: '深圳市', value: '深圳市', children: [
        { label: '罗湖区', value: '罗湖区' },
        { label: '福田区', value: '福田区' },
        { label: '南山区', value: '南山区' },
        { label: '宝安区', value: '宝安区' },
        { label: '龙岗区', value: '龙岗区' },
        { label: '盐田区', value: '盐田区' },
        { label: '龙华区', value: '龙华区' },
        { label: '坪山区', value: '坪山区' },
        { label: '光明区', value: '光明区' },
      ] },
      { label: '珠海市', value: '珠海市', children: [
        { label: '香洲区', value: '香洲区' },
        { label: '斗门区', value: '斗门区' },
        { label: '金湾区', value: '金湾区' },
      ] },
      { label: '东莞市', value: '东莞市', children: [
        { label: '莞城街道', value: '莞城街道' },
        { label: '南城街道', value: '南城街道' },
        { label: '东城街道', value: '东城街道' },
        { label: '万江街道', value: '万江街道' },
        { label: '虎门镇', value: '虎门镇' },
        { label: '长安镇', value: '长安镇' },
      ] },
      { label: '佛山市', value: '佛山市', children: [
        { label: '禅城区', value: '禅城区' },
        { label: '南海区', value: '南海区' },
        { label: '顺德区', value: '顺德区' },
        { label: '三水区', value: '三水区' },
        { label: '高明区', value: '高明区' },
      ] },
      { label: '惠州市', value: '惠州市', children: [
        { label: '惠城区', value: '惠城区' },
        { label: '惠阳区', value: '惠阳区' },
        { label: '博罗县', value: '博罗县' },
        { label: '惠东县', value: '惠东县' },
        { label: '龙门县', value: '龙门县' },
      ] },
    ]
  },
  {
    label: '浙江省', value: '浙江省', children: [
      { label: '杭州市', value: '杭州市', children: [
        { label: '上城区', value: '上城区' },
        { label: '拱墅区', value: '拱墅区' },
        { label: '西湖区', value: '西湖区' },
        { label: '滨江区', value: '滨江区' },
        { label: '萧山区', value: '萧山区' },
        { label: '余杭区', value: '余杭区' },
        { label: '临平区', value: '临平区' },
        { label: '钱塘区', value: '钱塘区' },
        { label: '富阳区', value: '富阳区' },
        { label: '临安区', value: '临安区' },
      ] },
      { label: '宁波市', value: '宁波市', children: [
        { label: '海曙区', value: '海曙区' },
        { label: '江北区', value: '江北区' },
        { label: '北仑区', value: '北仑区' },
        { label: '镇海区', value: '镇海区' },
        { label: '鄞州区', value: '鄞州区' },
        { label: '奉化区', value: '奉化区' },
      ] },
      { label: '温州市', value: '温州市', children: [
        { label: '鹿城区', value: '鹿城区' },
        { label: '龙湾区', value: '龙湾区' },
        { label: '瓯海区', value: '瓯海区' },
        { label: '洞头区', value: '洞头区' },
        { label: '瑞安市', value: '瑞安市' },
        { label: '乐清市', value: '乐清市' },
      ] },
      { label: '嘉兴市', value: '嘉兴市', children: [
        { label: '南湖区', value: '南湖区' },
        { label: '秀洲区', value: '秀洲区' },
        { label: '海宁市', value: '海宁市' },
        { label: '平湖市', value: '平湖市' },
        { label: '桐乡市', value: '桐乡市' },
      ] },
    ]
  },
  {
    label: '江苏省', value: '江苏省', children: [
      { label: '南京市', value: '南京市', children: [
        { label: '玄武区', value: '玄武区' },
        { label: '秦淮区', value: '秦淮区' },
        { label: '建邺区', value: '建邺区' },
        { label: '鼓楼区', value: '鼓楼区' },
        { label: '浦口区', value: '浦口区' },
        { label: '栖霞区', value: '栖霞区' },
        { label: '雨花台区', value: '雨花台区' },
        { label: '江宁区', value: '江宁区' },
      ] },
      { label: '苏州市', value: '苏州市', children: [
        { label: '虎丘区', value: '虎丘区' },
        { label: '吴中区', value: '吴中区' },
        { label: '相城区', value: '相城区' },
        { label: '姑苏区', value: '姑苏区' },
        { label: '吴江区', value: '吴江区' },
        { label: '昆山市', value: '昆山市' },
        { label: '常熟市', value: '常熟市' },
        { label: '张家港市', value: '张家港市' },
      ] },
      { label: '无锡市', value: '无锡市', children: [
        { label: '锡山区', value: '锡山区' },
        { label: '惠山区', value: '惠山区' },
        { label: '滨湖区', value: '滨湖区' },
        { label: '梁溪区', value: '梁溪区' },
        { label: '新吴区', value: '新吴区' },
        { label: '江阴市', value: '江阴市' },
        { label: '宜兴市', value: '宜兴市' },
      ] },
      { label: '常州市', value: '常州市', children: [
        { label: '天宁区', value: '天宁区' },
        { label: '钟楼区', value: '钟楼区' },
        { label: '新北区', value: '新北区' },
        { label: '武进区', value: '武进区' },
      ] },
    ]
  },
  {
    label: '四川省', value: '四川省', children: [
      { label: '成都市', value: '成都市', children: [
        { label: '锦江区', value: '锦江区' },
        { label: '青羊区', value: '青羊区' },
        { label: '金牛区', value: '金牛区' },
        { label: '武侯区', value: '武侯区' },
        { label: '成华区', value: '成华区' },
        { label: '龙泉驿区', value: '龙泉驿区' },
        { label: '青白江区', value: '青白江区' },
        { label: '新都区', value: '新都区' },
        { label: '温江区', value: '温江区' },
        { label: '双流区', value: '双流区' },
        { label: '郫都区', value: '郫都区' },
      ] },
      { label: '绵阳市', value: '绵阳市', children: [
        { label: '涪城区', value: '涪城区' },
        { label: '游仙区', value: '游仙区' },
        { label: '安州区', value: '安州区' },
      ] },
      { label: '德阳市', value: '德阳市', children: [
        { label: '旌阳区', value: '旌阳区' },
        { label: '罗江区', value: '罗江区' },
        { label: '广汉市', value: '广汉市' },
      ] },
    ]
  },
  {
    label: '湖北省', value: '湖北省', children: [
      { label: '武汉市', value: '武汉市', children: [
        { label: '江岸区', value: '江岸区' },
        { label: '江汉区', value: '江汉区' },
        { label: '硚口区', value: '硚口区' },
        { label: '汉阳区', value: '汉阳区' },
        { label: '武昌区', value: '武昌区' },
        { label: '青山区', value: '青山区' },
        { label: '洪山区', value: '洪山区' },
        { label: '东西湖区', value: '东西湖区' },
      ] },
      { label: '宜昌市', value: '宜昌市', children: [
        { label: '西陵区', value: '西陵区' },
        { label: '伍家岗区', value: '伍家岗区' },
        { label: '点军区', value: '点军区' },
        { label: '猇亭区', value: '猇亭区' },
      ] },
    ]
  },
  {
    label: '湖南省', value: '湖南省', children: [
      { label: '长沙市', value: '长沙市', children: [
        { label: '芙蓉区', value: '芙蓉区' },
        { label: '天心区', value: '天心区' },
        { label: '岳麓区', value: '岳麓区' },
        { label: '开福区', value: '开福区' },
        { label: '雨花区', value: '雨花区' },
        { label: '望城区', value: '望城区' },
        { label: '长沙县', value: '长沙县' },
      ] },
      { label: '株洲市', value: '株洲市', children: [
        { label: '天元区', value: '天元区' },
        { label: '荷塘区', value: '荷塘区' },
        { label: '芦淞区', value: '芦淞区' },
        { label: '石峰区', value: '石峰区' },
      ] },
    ]
  },
  {
    label: '山东省', value: '山东省', children: [
      { label: '济南市', value: '济南市', children: [
        { label: '历下区', value: '历下区' },
        { label: '市中区', value: '市中区' },
        { label: '槐荫区', value: '槐荫区' },
        { label: '天桥区', value: '天桥区' },
        { label: '历城区', value: '历城区' },
        { label: '长清区', value: '长清区' },
      ] },
      { label: '青岛市', value: '青岛市', children: [
        { label: '市南区', value: '市南区' },
        { label: '市北区', value: '市北区' },
        { label: '黄岛区', value: '黄岛区' },
        { label: '崂山区', value: '崂山区' },
        { label: '李沧区', value: '李沧区' },
        { label: '城阳区', value: '城阳区' },
      ] },
      { label: '烟台市', value: '烟台市', children: [
        { label: '芝罘区', value: '芝罘区' },
        { label: '福山区', value: '福山区' },
        { label: '牟平区', value: '牟平区' },
        { label: '莱山区', value: '莱山区' },
      ] },
    ]
  },
  {
    label: '河南省', value: '河南省', children: [
      { label: '郑州市', value: '郑州市', children: [
        { label: '中原区', value: '中原区' },
        { label: '二七区', value: '二七区' },
        { label: '管城回族区', value: '管城回族区' },
        { label: '金水区', value: '金水区' },
        { label: '上街区', value: '上街区' },
        { label: '惠济区', value: '惠济区' },
      ] },
      { label: '洛阳市', value: '洛阳市', children: [
        { label: '老城区', value: '老城区' },
        { label: '西工区', value: '西工区' },
        { label: '瀍河回族区', value: '瀍河回族区' },
        { label: '涧西区', value: '涧西区' },
        { label: '洛龙区', value: '洛龙区' },
      ] },
    ]
  },
  {
    label: '福建省', value: '福建省', children: [
      { label: '福州市', value: '福州市', children: [
        { label: '鼓楼区', value: '鼓楼区' },
        { label: '台江区', value: '台江区' },
        { label: '仓山区', value: '仓山区' },
        { label: '马尾区', value: '马尾区' },
        { label: '晋安区', value: '晋安区' },
        { label: '长乐区', value: '长乐区' },
      ] },
      { label: '厦门市', value: '厦门市', children: [
        { label: '思明区', value: '思明区' },
        { label: '海沧区', value: '海沧区' },
        { label: '湖里区', value: '湖里区' },
        { label: '集美区', value: '集美区' },
        { label: '同安区', value: '同安区' },
        { label: '翔安区', value: '翔安区' },
      ] },
      { label: '泉州市', value: '泉州市', children: [
        { label: '鲤城区', value: '鲤城区' },
        { label: '丰泽区', value: '丰泽区' },
        { label: '洛江区', value: '洛江区' },
        { label: '泉港区', value: '泉港区' },
        { label: '晋江市', value: '晋江市' },
        { label: '石狮市', value: '石狮市' },
      ] },
    ]
  },
  {
    label: '河北省', value: '河北省', children: [
      { label: '石家庄市', value: '石家庄市', children: [
        { label: '长安区', value: '长安区' },
        { label: '桥西区', value: '桥西区' },
        { label: '新华区', value: '新华区' },
        { label: '裕华区', value: '裕华区' },
      ] },
      { label: '唐山市', value: '唐山市', children: [
        { label: '路南区', value: '路南区' },
        { label: '路北区', value: '路北区' },
        { label: '古冶区', value: '古冶区' },
        { label: '开平区', value: '开平区' },
      ] },
      { label: '廊坊市', value: '廊坊市', children: [
        { label: '安次区', value: '安次区' },
        { label: '广阳区', value: '广阳区' },
        { label: '霸州市', value: '霸州市' },
        { label: '三河市', value: '三河市' },
      ] },
    ]
  },
  {
    label: '安徽省', value: '安徽省', children: [
      { label: '合肥市', value: '合肥市', children: [
        { label: '瑶海区', value: '瑶海区' },
        { label: '庐阳区', value: '庐阳区' },
        { label: '蜀山区', value: '蜀山区' },
        { label: '包河区', value: '包河区' },
      ] },
      { label: '芜湖市', value: '芜湖市', children: [
        { label: '镜湖区', value: '镜湖区' },
        { label: '弋江区', value: '弋江区' },
        { label: '鸠江区', value: '鸠江区' },
      ] },
    ]
  },
  {
    label: '江西省', value: '江西省', children: [
      { label: '南昌市', value: '南昌市', children: [
        { label: '东湖区', value: '东湖区' },
        { label: '西湖区', value: '西湖区' },
        { label: '青云谱区', value: '青云谱区' },
        { label: '青山湖区', value: '青山湖区' },
        { label: '新建区', value: '新建区' },
        { label: '红谷滩区', value: '红谷滩区' },
      ] },
    ]
  },
  {
    label: '辽宁省', value: '辽宁省', children: [
      { label: '沈阳市', value: '沈阳市', children: [
        { label: '和平区', value: '和平区' },
        { label: '沈河区', value: '沈河区' },
        { label: '大东区', value: '大东区' },
        { label: '皇姑区', value: '皇姑区' },
        { label: '铁西区', value: '铁西区' },
        { label: '浑南区', value: '浑南区' },
      ] },
      { label: '大连市', value: '大连市', children: [
        { label: '中山区', value: '中山区' },
        { label: '西岗区', value: '西岗区' },
        { label: '沙河口区', value: '沙河口区' },
        { label: '甘井子区', value: '甘井子区' },
        { label: '旅顺口区', value: '旅顺口区' },
        { label: '金州区', value: '金州区' },
      ] },
    ]
  },
  {
    label: '陕西省', value: '陕西省', children: [
      { label: '西安市', value: '西安市', children: [
        { label: '新城区', value: '新城区' },
        { label: '碑林区', value: '碑林区' },
        { label: '莲湖区', value: '莲湖区' },
        { label: '灞桥区', value: '灞桥区' },
        { label: '未央区', value: '未央区' },
        { label: '雁塔区', value: '雁塔区' },
        { label: '阎良区', value: '阎良区' },
        { label: '长安区', value: '长安区' },
      ] },
    ]
  },
  {
    label: '广西壮族自治区', value: '广西壮族自治区', children: [
      { label: '南宁市', value: '南宁市', children: [
        { label: '兴宁区', value: '兴宁区' },
        { label: '青秀区', value: '青秀区' },
        { label: '江南区', value: '江南区' },
        { label: '西乡塘区', value: '西乡塘区' },
        { label: '良庆区', value: '良庆区' },
        { label: '邕宁区', value: '邕宁区' },
      ] },
      { label: '桂林市', value: '桂林市', children: [
        { label: '秀峰区', value: '秀峰区' },
        { label: '叠彩区', value: '叠彩区' },
        { label: '象山区', value: '象山区' },
        { label: '七星区', value: '七星区' },
        { label: '雁山区', value: '雁山区' },
      ] },
    ]
  },
  {
    label: '云南省', value: '云南省', children: [
      { label: '昆明市', value: '昆明市', children: [
        { label: '五华区', value: '五华区' },
        { label: '盘龙区', value: '盘龙区' },
        { label: '官渡区', value: '官渡区' },
        { label: '西山区', value: '西山区' },
        { label: '呈贡区', value: '呈贡区' },
      ] },
    ]
  },
  {
    label: '贵州省', value: '贵州省', children: [
      { label: '贵阳市', value: '贵阳市', children: [
        { label: '南明区', value: '南明区' },
        { label: '云岩区', value: '云岩区' },
        { label: '花溪区', value: '花溪区' },
        { label: '乌当区', value: '乌当区' },
        { label: '白云区', value: '白云区' },
        { label: '观山湖区', value: '观山湖区' },
      ] },
    ]
  },
  {
    label: '黑龙江省', value: '黑龙江省', children: [
      { label: '哈尔滨市', value: '哈尔滨市', children: [
        { label: '道里区', value: '道里区' },
        { label: '南岗区', value: '南岗区' },
        { label: '道外区', value: '道外区' },
        { label: '平房区', value: '平房区' },
        { label: '松北区', value: '松北区' },
        { label: '香坊区', value: '香坊区' },
      ] },
    ]
  },
  {
    label: '吉林省', value: '吉林省', children: [
      { label: '长春市', value: '长春市', children: [
        { label: '南关区', value: '南关区' },
        { label: '宽城区', value: '宽城区' },
        { label: '朝阳区', value: '朝阳区' },
        { label: '二道区', value: '二道区' },
        { label: '绿园区', value: '绿园区' },
      ] },
    ]
  },
  {
    label: '山西省', value: '山西省', children: [
      { label: '太原市', value: '太原市', children: [
        { label: '小店区', value: '小店区' },
        { label: '迎泽区', value: '迎泽区' },
        { label: '杏花岭区', value: '杏花岭区' },
        { label: '尖草坪区', value: '尖草坪区' },
        { label: '万柏林区', value: '万柏林区' },
        { label: '晋源区', value: '晋源区' },
      ] },
    ]
  },
  {
    label: '甘肃省', value: '甘肃省', children: [
      { label: '兰州市', value: '兰州市', children: [
        { label: '城关区', value: '城关区' },
        { label: '七里河区', value: '七里河区' },
        { label: '西固区', value: '西固区' },
        { label: '安宁区', value: '安宁区' },
        { label: '红古区', value: '红古区' },
      ] },
    ]
  },
  {
    label: '内蒙古自治区', value: '内蒙古自治区', children: [
      { label: '呼和浩特市', value: '呼和浩特市', children: [
        { label: '新城区', value: '新城区' },
        { label: '回民区', value: '回民区' },
        { label: '玉泉区', value: '玉泉区' },
        { label: '赛罕区', value: '赛罕区' },
      ] },
      { label: '包头市', value: '包头市', children: [
        { label: '东河区', value: '东河区' },
        { label: '昆都仑区', value: '昆都仑区' },
        { label: '青山区', value: '青山区' },
      ] },
    ]
  },
  {
    label: '海南省', value: '海南省', children: [
      { label: '海口市', value: '海口市', children: [
        { label: '秀英区', value: '秀英区' },
        { label: '龙华区', value: '龙华区' },
        { label: '琼山区', value: '琼山区' },
        { label: '美兰区', value: '美兰区' },
      ] },
      { label: '三亚市', value: '三亚市', children: [
        { label: '海棠区', value: '海棠区' },
        { label: '吉阳区', value: '吉阳区' },
        { label: '天涯区', value: '天涯区' },
        { label: '崖州区', value: '崖州区' },
      ] },
    ]
  },
  {
    label: '宁夏回族自治区', value: '宁夏回族自治区', children: [
      { label: '银川市', value: '银川市', children: [
        { label: '兴庆区', value: '兴庆区' },
        { label: '西夏区', value: '西夏区' },
        { label: '金凤区', value: '金凤区' },
      ] },
    ]
  },
  {
    label: '新疆维吾尔自治区', value: '新疆维吾尔自治区', children: [
      { label: '乌鲁木齐市', value: '乌鲁木齐市', children: [
        { label: '天山区', value: '天山区' },
        { label: '沙依巴克区', value: '沙依巴克区' },
        { label: '新市区', value: '新市区' },
        { label: '水磨沟区', value: '水磨沟区' },
        { label: '头屯河区', value: '头屯河区' },
      ] },
    ]
  },
  {
    label: '西藏自治区', value: '西藏自治区', children: [
      { label: '拉萨市', value: '拉萨市', children: [
        { label: '城关区', value: '城关区' },
        { label: '堆龙德庆区', value: '堆龙德庆区' },
        { label: '达孜区', value: '达孜区' },
      ] },
    ]
  },
  {
    label: '青海省', value: '青海省', children: [
      { label: '西宁市', value: '西宁市', children: [
        { label: '城东区', value: '城东区' },
        { label: '城中区', value: '城中区' },
        { label: '城西区', value: '城西区' },
        { label: '城北区', value: '城北区' },
      ] },
    ]
  },
]

// 获取城市列表（根据省份）
export function getCities(province: string): RegionOption[] {
  const prov = regionData.find(p => p.value === province)
  return prov?.children || []
}

// 获取区县列表（根据省份和城市）
export function getDistricts(province: string, city: string): RegionOption[] {
  const cities = getCities(province)
  const c = cities.find(c => c.value === city)
  return c?.children || []
}

var ViewMode = true; // xem theo vacxin
var DsVacxin = null;
var IsLichSuCapNhatLoaded = false;
var IsDinhDuongLoaded = false;
var IsDichVuDangKyLoaded = false;
InitFormDetail();

function InitFormDetail() {
    var slXa = $("#slXaDetail").select2({});
    var slThonAp = $("#slThonApDetail").select2({});
    var slDanToc = $("#slDanTocDetail").select2({ width: "100%" });

    var slXaDangKy = $("#slXaDangKyDetail").select2();
    var slThonApDangKy = $("#slThonApDangKyDetail").select2();

    var slNhomMau = $("#slNhomMauDetail").select2({});
    var slTonGiao = $("#slTonGiaoDetail").select2({});
    var slQuocGia = $("#slQuocTichDetail").select2({});

    var slDVHCNoiSinh = $("#slDVHCNsDetail").select2();
    var slDVHCDkks = $("#slDVHCDkksDetail").select2();
    var slDVHCQq = $("#slDVHCQqDetail").select2();

    var slThonApNs = $("#slThonApNsDetail").select2({});
    var slThonApDkks = $("#slThonApDkksDetail").select2({});
    var slThonApQq = $("#slThonApQqDetail").select2({});
    var slTinhTrangCd = $("#slTinhTrangCdDetail").select2({});

    var tggNgungTheoDoi = $('#tggNgungTheoDoi');
    var hfNgungTheoDoi = $("#hfNgungTheoDoi").val() == 1 ? true : false;

    tggNgungTheoDoi.toggles({
        drag: true,
        click: true,
        text: {
            on: GlobalResources.DT_BAT,
            off: GlobalResources.DT_TAT
        },
        on: hfNgungTheoDoi,
        animate: 250,
        easing: 'swing',
        checkbox: null,
        clicker: null,
        width: 60,
        height: 25,
        type: 'compact',
        disabled: true
    });

    var tggBVUVSS = $('#tggBVUVSS');
    var hfBVUVSS = $("#hfBVUVSS").val() == 1 ? true : false;

    //nếu "Được bảo vệ UVSS" = true > hiển thị (add hoangnv 28/08/2019)
    if (hfBVUVSS == true) {
        $("#divBvuvss").show();
    } else { //không -> ẩn
        $("#divBvuvss").hide();
    }

    tggBVUVSS.toggles({
        drag: true,
        click: true,
        text: {
            on: GlobalResources.DT_CO,
            off: GlobalResources.DT_KHONG
        },
        on: hfBVUVSS,
        animate: 250,
        easing: 'swing',
        checkbox: null,
        clicker: null,
        width: 80,
        height: 25,
        type: 'compact',
        disabled: true
    });

    var tggRaSoatTienSuMuiTiem = $('#tggRaSoatTienSuMuiTiem');
    var hfRaSoatTienSuMuiTiem = $("#hfRaSoatTienSuMuiTiem").val() == 1 ? true : false;

    tggRaSoatTienSuMuiTiem.toggles({
        drag: true,
        click: true,
        text: {
            on: GlobalResources.DT_CO,
            off: GlobalResources.DT_KHONG
        },
        on: hfRaSoatTienSuMuiTiem,
        animate: 250,
        easing: 'swing',
        checkbox: null,
        clicker: null,
        width: 60,
        height: 25,
        type: 'compact',
        disabled: true
    });

    var tggBaoVeDuLieu = $('#tggBaoVeDuLieu');
    var hfBaoVeDuLieu = $("#hfBaoVeDuLieu").val() == 1 ? true : false;

    tggBaoVeDuLieu.toggles({
        drag: true,
        click: true,
        text: {
            on: GlobalResources.DT_CO,
            off: GlobalResources.DT_KHONG
        },
        on: hfBaoVeDuLieu,
        animate: 250,
        easing: 'swing',
        checkbox: null,
        clicker: null,
        width: 60,
        height: 25,
        type: 'compact',
        disabled: true
    });

    var tggKbtvDetail = $('#tggKbtvDetail');
    var hfKbtvDetail = $("#hfKbtvDetail").val() == 1 ? true : false;

    tggKbtvDetail.toggles({
        drag: true,
        click: true,
        text: {
            on: GlobalResources.DT_CO,
            off: GlobalResources.DT_KHONG
        },
        on: hfKbtvDetail,
        animate: 250,
        easing: 'swing',
        checkbox: null,
        clicker: null,
        width: 60,
        height: 25,
        type: 'compact',
        disabled: true
    });

    $('#tggKbtvDetail')
        .css('pointer-events', 'none')
        .css('opacity', '0.6');

    $(".spinner").spinner("disable");
    $("#slGioiTinh_Detail").select2({ width: "100%", minimumResultsForSearch: -1 });

    var typeLoaiDoiTuong = $("#hfLoaiDoiTuong_Detail").val();
    if (typeLoaiDoiTuong == 1) {         // Tre em
        InitTreEmForm();
    } else if (typeLoaiDoiTuong == 2) {  // Phu nu
        InitPhuNuForm();
    } else if (typeLoaiDoiTuong == 3) {  // Khac
        InitDoiTuongKhacForm();
    }

    // che do xem phac do
    var tggViewMode = $('#tggViewMode');

    tggViewMode.toggles({
        drag: true,
        click: true,
        text: {
            on: GlobalResources.DT_VAC_XIN,
            off: GlobalResources.DT_KHANG_NGUYEN
        },
        on: true,
        animate: 250,
        easing: 'swing',
        checkbox: null,
        clicker: null,
        width: 100,
        height: 20,
        type: 'compact'
    });

    tggViewMode.on('toggle', function (e, active) {
        ViewMode = active;
        if (ViewMode) {
            XemTheoVacxin();
        } else {
            XemTheoKhangNguyen();
        }
        if (TiemChung$DoiTuong$Detail$DoiTuongModel != null)
            if (TiemChung$DoiTuong$Detail$DoiTuongModel.IS_ACTIVE == 1) {
                $("#btnThemMuiTiem").hide();
                $("#btnThemNhanh").hide();
                $('a[name="btnActionLST"]').hide();
            }
    });

    $("#btnThemMuiTiem").on("click", function (e) {
        ShowThemMuiTiemDialog();
    });

    $("#btnThemNhanh").on("click", function (e) {
        ThemNhanhMuiTiem();
    });

    $("#btnKiemTraTCDD").on("click", function (e) {
        KiemTraTCDD();
    });

    $("#btnThemSoDoDinhDuong").on("click", function (e) {
        ThemSoDoDinhDuong("INSERT");
    });

    $("#btnCanNang").on("click", function (e) {
        ShowModalBieuDoDinhDuongCanNang();
    });

    $("#btnChieuCao").on("click", function (e) {
        ShowModalBieuDoDinhDuongChieuCao();
    });

    $("#tbLichSuCapNhat").on("shown.bs.tab", function () {
        if (!IsLichSuCapNhatLoaded) {
            XemLichSuCapNhat();
        }
    });

    $("#tbThongTinDinhDuong").on("show.bs.tab", function () {
        XemThongTinDinhDuong();
    });

    $("#tbThongTinDichVu").on("show.bs.tab", function () {
        XemThongTinDichVu();
    });

    $("#btnDangKyDichVu").on("click", function () {
        ShowModalDangKyDichVu();
    });

    $("#btnLichSuTinNhan").on("click", function () {
        ShowModalLichSuTinNhan();
    });

    ShowActionButtonForDetail();
    KhoiTaoThongTinDiaChiDoiTuong();
    KhoiTaoThongTinDanToc_Detail();
    KhoiTaoDuLieuTonGiao_Detail();
    KhoiTaoDuLieuQuocTich_Detail();
}

function XemThongTinDinhDuong() {
    if (!IsDinhDuongLoaded) {
        Common.UI.BlockElement("#thongTinDinhDuong");
        $.ajax({
            url: "/TiemChung/DoiTuong/ThongTinDinhDuong",
            type: "GET",
            async: true,
            data: {
                DoiTuongId: TiemChung$DoiTuong$Detail$DoiTuongModel.DOI_TUONG_2_ID,
                PageNumber: 1,
                PageSize: 20
            },
            success: function (response) {
                IsDinhDuongLoaded = true;

                $("#divDinhDuong").empty();
                $("#divDinhDuong").append(response);
            },
            error: function (e) {
                jQuery.gritter.add({
                    title: GlobalResources.DT_LOI,
                    text: GlobalResources.DT_KHONG_THE_LAY_THONG_TIN_DINH_DUONG_DT,
                    class_name: 'growl-danger',
                    sticky: false,
                    time: 2000
                });
            },
            complete: function () {
                Common.UI.UnBlockElement("#thongTinDinhDuong");
            }
        });
    }
}

function XemThongTinDichVu() {
    if (!IsDichVuDangKyLoaded) {
        Common.UI.BlockElement("#thongTinDichVu");
        $.ajax({
            url: "/TiemChung/DoiTuong/DichVuDangKy",
            type: "GET",
            async: true,
            data: {
                DoiTuongId: TiemChung$DoiTuong$Detail$DoiTuongModel.DOI_TUONG_2_ID
            },
            success: function (response) {
                $("#divDanhSachDichVu").empty();
                $("#divDanhSachDichVu").append(response);
                IsDichVuDangKyLoaded = true;
                Common.UI.UnBlockElement("#thongTinDichVu");
            },
            error: function () {
                jQuery.gritter.add({
                    title: GlobalResources.DT_LOI,
                    text: GlobalResources.DT_KHONG_THE_LAY_DUOC_THONG_TIN_DVDK_DT,
                    class_name: 'growl-danger',
                    sticky: false,
                    time: 2000
                });
            }
        });
    }
}

function XemLichSuCapNhat() {

    Common.UI.BlockElement("#thongTinCapNhat");
    $.ajax({
        url: "/TiemChung/DoiTuong/LichSuCapNhat",
        type: "GET",
        async: true,
        data: { DoiTuongId: TiemChung$DoiTuong$Detail$DoiTuongModel.DOI_TUONG_2_ID },
        success: function (response) {
            $("#thongTinCapNhat").empty();
            $("#thongTinCapNhat").append(response);

            IsLichSuCapNhatLoaded = true;
            Common.UI.UnBlockElement("#thongTinCapNhat");
        },
        error: function (e) {
            jQuery.gritter.add({
                title: GlobalResources.DT_LOI,
                text: GlobalResources.DT_KHONG_THE_LAY_DUOC_THONG_TIN_CAP_NHAT_CUA_DT,
                class_name: 'growl-danger',
                sticky: false,
                time: 2000
            });
        }
    });
}

function KhoiTaoThongTinDiaChiDoiTuong() {

    // khoi tao dia chi ho khau
    if (TiemChung$DoiTuong$Detail$DoiTuongModel != null) {
        if (TiemChung$DoiTuong$Detail$DoiTuongModel.XA_ID != null) {
            ImCache.GetXaByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.XA_ID, function (xa) {
                var optionXaDetail = $("<option selected />");
                optionXaDetail.val(xa.XA_ID);
                optionXaDetail.text(xa.TEN_XA_TINH);

                $("#slXaDetail").append(optionXaDetail);
                $("#slXaDetail").select2({});
            });
        }

        // khoi tao dia chi dang ky        
        if (TiemChung$DoiTuong$Detail$DoiTuongModel.XA_DANG_KY_ID != null) {
            ImCache.GetXaByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.XA_DANG_KY_ID, function (xaDangKy) {
                var optionXaDangKyDetail = $("<option selected />");
                optionXaDangKyDetail.val(xaDangKy.XA_ID);
                optionXaDangKyDetail.text(xaDangKy.TEN_XA_TINH);

                $("#slXaDangKyDetail").append(optionXaDangKyDetail);
                $("#slXaDangKyDetail").select2({});
            });
        }

        if (TiemChung$DoiTuong$Detail$DoiTuongModel.NOI_SINH_XA_ID != null) {
            ImCache.GetXaByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.NOI_SINH_XA_ID, function (xaNoiKhaiSinh) {
                var option = $("<option selected />")
                    .val(xaNoiKhaiSinh.XA_ID)
                    .text(batBaoVeDuLieu ? "********" : xaNoiKhaiSinh.TEN_XA_TINH);

                $("#slDVHCNsDetail").append(option).select2({});
            });
        }

        if (TiemChung$DoiTuong$Detail$DoiTuongModel.NOI_KHAI_SINH_XA_ID != null) {
            ImCache.GetXaByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.NOI_KHAI_SINH_XA_ID, function (xaNoiKhaiSinh) {
                var option = $("<option selected />")
                    .val(xaNoiKhaiSinh.XA_ID)
                    .text(batBaoVeDuLieu ? "********" : xaNoiKhaiSinh.TEN_XA_TINH);

                $("#slDVHCDkksDetail").append(option).select2({});
            });
        }

        if (TiemChung$DoiTuong$Detail$DoiTuongModel.QUE_QUAN_XA_ID != null) {
            ImCache.GetXaByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.QUE_QUAN_XA_ID, function (xaQueQuan) {
                var option = $("<option selected />")
                    .val(xaQueQuan.XA_ID)
                    .text(batBaoVeDuLieu ? "********" : xaQueQuan.TEN_XA_TINH);

                $("#slDVHCQqDetail").append(option).select2({});
            });
        }
    }
}

function KhoiTaoThongTinDanToc_Detail() {
    if (TiemChung$DoiTuong$Detail$DoiTuongModel != null) {
        ImCache.GetDanTocByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.DAN_TOC_ID, function (danToc) {
            if (danToc != null) {
                var danTocOption = $("<option selected='selected' />")
                danTocOption.val(danToc.DAN_TOC_ID);
                danTocOption.text(danToc.TEN_DAN_TOC);

                $("#slDanTocDetail").append(danTocOption);
            }
        });
    }
}

function KhoiTaoDuLieuQuocTich_Detail() {
    if (TiemChung$DoiTuong$Detail$DoiTuongModel != null) {
        ImCache.GetQuocGiaByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.QUOC_GIA_ID, function (quocGia) {
            const $slQuocTich = $("#slQuocTichDetail");

            const option = $("<option>", { selected: true });
            if (quocGia) {
                option.val(quocGia.QUOC_GIA_ID)
                    .text(batBaoVeDuLieu ? "********" : quocGia.TEN_QUOC_GIA);
            } else if (batBaoVeDuLieu) {
                option.val(0).text("********");
            } else {
                return;
            }

            $slQuocTich.append(option);
        });
    }
}

function KhoiTaoDuLieuTonGiao_Detail() {

    if (TiemChung$DoiTuong$Detail$DoiTuongModel != null) {
        ImCache.GetTonGiaoByIdFromCache(TiemChung$DoiTuong$Detail$DoiTuongModel.TON_GIAO_ID, function (tonGiao) {
            const $slTonGiao = $("#slTonGiaoDetail");

            const option = $("<option>", { selected: true });
            if (tonGiao) {
                option.val(tonGiao.TON_GIAO_ID)
                    .text(batBaoVeDuLieu ? "********" : tonGiao.TEN_TON_GIAO);
            } else if (batBaoVeDuLieu) {
                option.val(0).text("********");
            } else {
                return;
            }

            $slTonGiao.append(option);
        });
    }
}

function KiemTraTCDD() {
    var OBJ_PARAMS_TINHTCDD = {
        'CAP_ID': 0,
        'TINH_ID': TiemChung$DoiTuong$Detail$DoiTuongModel.DOI_TUONG_2_ID,
        'HUYEN_ID': 0,
        'XA_ID': 0,
        'TU_NGAY': 0,
        'TOI_NGAY': 0
    };
    // khoa nut kiem tra TCDD
    $("#btnKiemTraTCDD").prop("disabled", true);
    $.ajax({
        type: "POST",
        url: '/TiemChung/DoiTuong/TinhTCDD',
        data: JSON.stringify(OBJ_PARAMS_TINHTCDD),
        headers: layGiaTriToken(),
        dataType: "json",
        async: false,
        traditional: true,
        contentType: "application/json; charset=utf-8",
        success: function (data) {
            if (data.Status == 1) {
                CommonJS.showSuccessMessage(data.Message);
            } else {
                CommonJS.showDangerMessage(data.Message);
            }
            $("#btnKiemTraTCDD").prop("disabled", false);
        },
        error: function () {
            CommonJS.showDangerMessage(data.Message);
            $("#btnKiemTraTCDD").prop("disabled", false);
        },
        complete: function () {
            CommonJS.showDangerMessage(data.Message);
        }
    });
}

function ThemNhanhMuiTiem() {
    // khoa nut them nhanh
    $("#btnThemNhanh").prop("disabled", true);

    var cosoId = $('#hfCosoId').val();
    var tenCoSo = $('#hfTenCoSo').val();

    // tạo một row mới
    var rowHtml = '<tr> \
                            <td></td> \
                            <td><select class="vacxin" name="Vacxin" id="qaVacxin"><option selected="selected">- '+ GlobalResources.DT_CHON_VAC_XIN + ' -</option></select><span name="KhangNguyen" class="sublabel"></span></td> \
                            <td><input id="SoMuiUV_ThemNhanh" type="text" class="form-control input-sm" value="" maxlength="2" style="display: none;"></td> \
                            <td><select class="form-control input-sm" name="TrangThai" id="qaTrangThai"><option value="2">' + GlobalResources.DT_DA_TIEM + '</option><option value="3">' + GlobalResources.DT_CHONG_CHI_DINH + '</option><option value="4">' + GlobalResources.KHT_KHONG_DONG_Y_TIEM + '</option><option value="5">' + GlobalResources.KHT_TAM_HOAN_TIEM_CHUNG + '</option></select></td> \
                            <td><div class="input-group input-group-sm" style="width: 100%"> \
                                    <input name="NgayTiem" type="text" class="form-control date-picker" placeholder="'+ GlobalResources.DT_NGAY_THANG_NAM + '" value="" id="qaNgayTiem"> \
                                    <span class="input-group-addon"><i class="glyphicon glyphicon-calendar"></i></span> \
                                </div> \
                            </td> \
                            <td><select name="CoSo" id="qaCoSo"><option value="'+ cosoId + '" selected="selected">' + tenCoSo + '</select></td> \
                            <td></td> \
                            <td style="text-align: center"> \
                                <a href="#" class="fa-hover mr10" title="' + GlobalResources.DT_LUU_MUI_TIEM + '" data-toggle="tooltip" onclick="LuuMuiTiem(this)"> \
                                    <i class="fa fa-save color-primary"></i> \
                                </a> \
                                <a href="#" class="fa-hover mr10" title="' + GlobalResources.DT_HUY + '" data-toggle="tooltip" onclick="HuyThemMuiTiem(this)"> \
                                    <i class="fa fa-reply color-primary"></i> \
                                </a> \
                            </td> \
                       </tr>';

    var rowObj = $(rowHtml);
    var colTrangThai = rowObj.find("select[name='TrangThai']").first();
    var colVacxin = rowObj.find("select[name='Vacxin']").first();
    var labelKhangNguyen = rowObj.find("span[name='KhangNguyen']").first();
    var colCoSo = rowObj.find("select[name='CoSo']").first();
    var colNgayTiem = rowObj.find("input[name='NgayTiem']").first();

    jQuery.datetimepicker.setLocale(GlobalResources.DATE_TIME_PICKER_LANGUAGE);
    $(colNgayTiem).datetimepicker({
        format: 'd/m/Y',
        step: 1,
        i18n: {
            vi: {
                months: [GlobalResources.THANG_MOT, GlobalResources.THANG_HAI, GlobalResources.THANG_BA, GlobalResources.THANG_TU,
                GlobalResources.THANG_NAM, GlobalResources.THANG_SAU, GlobalResources.THANG_BAY, GlobalResources.THANG_TAM, GlobalResources.THANG_CHIN,
                GlobalResources.THANG_MUOI, GlobalResources.THANG_MUOI_MOT, GlobalResources.THANG_MUOI_HAI],
                dayOfWeek: [GlobalResources.CHU_NHAT, GlobalResources.THU_HAI, GlobalResources.THU_BA, GlobalResources.THU_TU, GlobalResources.THU_NAM, GlobalResources.THU_SAU, GlobalResources.THU_BAY]
            }
        },
        mask: '99/99/9999',
        timepicker: false
    });
    $(colNgayTiem).mask("99/99/9999");

    $(colVacxin).select2({
        placeholder: GlobalResources.DT_CHON_VAC_XIN,
        allowClear: true,
        width: '100%'
    });
    $(colTrangThai).select2({
        width: '100%'
    });

    $(colVacxin).on("change", function () {
        if ($(colVacxin).select2("val") != null) {
            var selectedVacxinItem = DsVacxin.filter(function (vacxinItem) {
                if (vacxinItem.VACXIN_ID == $(colVacxin).select2("val"))
                    return vacxinItem;
            })[0];
            if (selectedVacxinItem.KHANG_NGUYEN == GlobalResources.DT_UON_VAN && $("#hfLoaiDoiTuong_Detail").val() == 2) {
                $("#SoMuiUV_ThemNhanh").show();
            } else {
                $("#SoMuiUV_ThemNhanh").hide();
                $("#SoMuiUV_ThemNhanh").val('');
            }
            $(labelKhangNguyen).text(selectedVacxinItem.KHANG_NGUYEN);

            if (selectedVacxinItem.VACXIN_ID == 4)// opv -> da uong
            {
                (colTrangThai).find('option:first-child').text(GlobalResources.DT_DA_UONG);
            }
            else {
                (colTrangThai).find('option:first-child').text(GlobalResources.KHT_DA_TIEM);
            }
            $(colTrangThai).select2({
                width: '100%'
            });

        } else {
            $(labelKhangNguyen).text('');
            $("#SoMuiUV_ThemNhanh").hide();
            $("#SoMuiUV_ThemNhanh").val('');
        }
    });

    // kiem tra xem vacxin da duoc tai chua
    if (DsVacxin == null) {
        ImCache.GetDsVacxinKhongCovidFromCache(function (dsVacxinCached) {
            DsVacxin = dsVacxinCached;
            KhoiTaoHopChonVacxin();
        });
    } else {
        KhoiTaoHopChonVacxin();
    }

    function KhoiTaoHopChonVacxin() {
        var dsVacxinMoRong = DsVacxin.filter(function (vacxin) {
            return vacxin.TCMR == 1;
        }).sort(function (a, b) {
            return (a.VACXIN_ID - b.VACXIN_ID);
        });

        var dsVacxinDichVu = DsVacxin.filter(function (vacxin) {
            return vacxin.TCMR == 2;
        });

        var optGroupMoRong = $("<optgroup label='" + GlobalResources.DT_VAC_XIN_MO_RONG + "'></optgroup>");
        $.each(dsVacxinMoRong, function (index, item) {
            var option = $("<option />");
            $(option).val(item.VACXIN_ID);
            $(option).text(item.TEN_VACXIN);

            optGroupMoRong.append(option);
        });

        var optGroupDichVu = $("<optgroup label='" + GlobalResources.DT_VAC_XIN_DICH_VU + "'></optgroup>");
        $.each(dsVacxinDichVu, function (index, item) {
            var option = $("<option />");
            $(option).val(item.VACXIN_ID);
            $(option).text(item.TEN_VACXIN);

            optGroupDichVu.append(option);
        });

        $(colVacxin).append(optGroupMoRong);
        $(colVacxin).append(optGroupDichVu);
        $(colVacxin).select2({
            placeholder: GlobalResources.DT_CHON_VAC_XIN,
            allowClear: true,
            width: '100%'
        });
    }

    $(colCoSo).select2({
        ajax: {
            url: "/DonViHanhChinh/TimKiemCoSo",
            dataType: 'json',
            delay: 250,
            async: true,
            data: function (params) {
                return {
                    keyword: params.term,
                    page: params.page || 1
                };
            },
            processResults: function (data, params) {

                params.page = params.page || 1;
                return {
                    results: $.map(data, function (item) {
                        return {
                            id: item.CO_SO_ID,
                            text: item.TEN_CO_SO,
                            TEN_CO_SO: item.TEN_CO_SO,
                            TEN_XA: item.TEN_XA,
                            TEN_HUYEN: item.TEN_HUYEN,
                            TEN_TINH: item.TEN_TINH
                        }
                    }),
                    pagination: {
                        more: (data.length == 20 ? true : false)
                    }
                };
            },

            cache: true
        },
        "language": "vi",
        templateResult: function (repo) {
            if (repo.loading) return repo.text;

            var markup = "<div class='select2-result-repository clearfix'>\
                                <div class='select2-result-repository__meta'>\
                                    <div class='select2-result-repository__title'>"
                + repo.TEN_CO_SO + "</div>"
                + "<div class='select2-result-repository__description'>"
                + (repo.TEN_XA != null ? repo.TEN_XA + ", " : "")
                + (repo.TEN_HUYEN != null ? repo.TEN_HUYEN + ", " : "")
                + (repo.TEN_TINH != null ? repo.TEN_TINH : "")
                + "</div><div><div>";

            return markup;
        },
        templateSelection: function (repo) {
            return repo.text || repo.TEN_CO_SO;
        },
        escapeMarkup: function (markup) { return markup; },
        minimumInputLength: 3,
        placeholder: GlobalResources.DT_CHON_CO_SO_TIEM_CHUNG,
        allowClear: true,
        width: '100%'
    });

    $("#tblVacxin > tbody").prepend(rowObj);
}

var ForceSave = 0;

function LuuMuiTiem(obj) {
    var selectedRow = $(obj).parent().parent();
    var data = {
        VACXIN_ID: $("#qaVacxin").val(),
        TRANG_THAI_MUI_TIEM: $("#qaTrangThai").val(),
        SO_MUI_UV: $("#SoMuiUV_ThemNhanh").val(),
        NGAY_TIEM: $("#qaNgayTiem").val(),
        CO_SO_ID: $("#qaCoSo").val(),
        DOI_TUONG_ID: $("#hfDoiTuongId_Detail").val(),
        FORCE_SAVE: ForceSave
    };

    var validateData = (function (data) {
        // Kiểm tra tính hợp lệ của vắc xin
        if (data.VACXIN_ID == null || isNaN(data.VACXIN_ID)) {
            return {
                isValid: false,
                Message: GlobalResources.DT_BAN_PHAI_CHON_VAC_XIN
            }
        }

        if ($("#SoMuiUV_ThemNhanh").is(":visible")) {
            if (data.SO_MUI_UV != null && data.SO_MUI_UV != '') {
                if (isNaN(data.SO_MUI_UV)) {
                    return {
                        isValid: false,
                        Message: GlobalResources.DT_SO_MUI_UV_PHAI_LA_SO
                    }
                } else {
                    if (data.SO_MUI_UV <= 0) {
                        return {
                            isValid: false,
                            Message: GlobalResources.DT_SO_MUI_UV_PHAI_LON_HON_0
                        }
                    }
                }
            } else {
                return {
                    isValid: false,
                    Message: GlobalResources.DT_BAN_PHAI_NHAP_GIA_TRI_SO_MUI_UV
                }
            }

        }

        // Kiểm tra tính hợp lệ ngày tiêm
        var ngayTiemDateTime = moment(data.NGAY_TIEM, "DD/MM/YYYY");
        var ngayHienTai = moment($("#hfNgayHienTai").val(), "DD/MM/YYYY");

        if (ngayTiemDateTime > ngayHienTai) {
            return {
                isValid: false,
                Message: GlobalResources.DT_NGAY_TIEM_KHONG_LON_HON_NGAY_HIEN_TAI
            }
        }

        var ngaySinh = moment($("#txtNgaySinh").val(), "DD/MM/YYYY");
        if (ngayTiemDateTime < ngaySinh) {
            return {
                isValid: false,
                //Message: GlobalResource.DT_NGAY_TIEM_KHONG_DUOC_NHO_HON_NGAY_SINH + " (" + $("#txtNgaySinh").val() + ")"
                Message: "Ngày tiêm không được nhỏ hơn ngày sinh" + " (" + $("#txtNgaySinh").val() + ")"
            }
        }

        return {
            isValid: true,
            Message: GlobalResources.DT_THONG_TIN_HOP_LE
        }
    })(data);

    if (!validateData.isValid) {
        jQuery.gritter.add({
            title: GlobalResources.DT_THONG_BAO,
            text: validateData.Message,
            class_name: 'growl-danger',
            sticky: false,
            time: 2000
        });
    } else {
        $.ajax({
            url: "/TiemChung/DoiTuong/ThemNhanhMuiTiem",
            data: data,
            type: "POST",
            dataType: "json",
            beforeSend: function () {
                Common.UI.BlockElement(".contentpanel");
            },
            success: function (response) {
                if (response.Status == 1) {
                    // load lại danh sách mũi tiêm
                    RefreshLichSuTiemChung();
                    $("#btnThemNhanh").prop("disabled", false);
                } else if (response.Status == 0 || response.Status == 2 || response.Status == 3) {
                    var divError = $("<div></div>");
                    $.each(response.Message, function (index, item) {
                        var liError = $("<div></div>");
                        liError.text(item.Value);
                        divError.append(liError);
                    });

                    // Thong bao loi tu server
                    jQuery.gritter.add({
                        text: divError.html(),
                        class_name: 'growl-danger',
                        sticky: false,
                        timeout: 2000
                    });
                } else if (response.Status == 4) // warning
                {
                    bootbox.confirm({
                        size: 'small',
                        message: response.Message + '. ' + GlobalResources.DT_BAN_CO_MUON_TIEP_TUC_THEM_MUI_TIEM_NAY_KHONG + '?',
                        title: GlobalResources.DT_XAC_NHAN_THEM_MUI_TIEM,
                        backdrop: false,
                        callback: function (result) {
                            if (result != null) {
                                if (result) {
                                    ForceSave = 1;
                                    LuuMuiTiem(obj); // 1 -- force save, 0 -- default
                                    ForceSave = 0;
                                }
                            }
                        },
                        buttons: {
                            confirm: {
                                label: GlobalResources.DT_TIEP_TUC,
                                className: 'btn-warning btn-sm'
                            },
                            cancel: {
                                label: GlobalResources.DT_DONG,
                                className: 'btn-default btn-sm'
                            }
                        }
                    });
                }
            },
            error: function (e) {
                jQuery.gritter.add({
                    title: GlobalResources.DT_THONG_BAO,
                    text: GlobalResources.DT_CO_LOI_XAY_RA_TRONG_XU_LY_THU_LAI,
                    class_name: 'growl-danger',
                    sticky: false,
                    timeout: 2000
                });
                $("#btnThemNhanh").prop("disabled", false);
                Common.UI.UnBlockElement(".contentpanel");
            },
            complete: function () {
                Common.UI.UnBlockElement(".contentpanel");
            }
        });
    }
}

function HuyThemMuiTiem(obj) {
    // Xóa row đang thao tác
    var newRow = $(obj).parent().parent().remove();

    // Mở lại nút Thêm mũi tiêm
    $("#btnThemNhanh").prop("disabled", false);
}

function XemTheoVacxin() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();
    $("#btnThemNhanh").show();
    $("#btnThemNhanh").prop("disabled", false);
    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/LichSuTiemChung",
        data: { doiTuongId: doiTuongId, viewMode: true },
        type: "GET",
        async: true,
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            if (response.Status == 0) {
                jQuery.gritter.add({
                    title: GlobalResources.DT_THONG_BAO,
                    message: response.Message,
                    sticky: false,
                    timeout: '',
                    class_name: 'growl-danger'
                });
            } else {
                $("#lichSuTiem").empty();
                $("#lichSuTiem").html(response);
                if (TiemChung$DoiTuong$Detail$DoiTuongModel != null)
                    if (TiemChung$DoiTuong$Detail$DoiTuongModel.IS_ACTIVE == 1) {
                        $('a[name="btnActionLST"]').hide();
                    }
            }
        },
        error: function (error) {
            Common.UI.UnBlockElement($("#detailPanel"));
            jQuery.gritter.add({
                title: GlobalResources.DT_THONG_BAO,
                message: error.message,
                sticky: false,
                timeout: '',
                class_name: 'growl-danger'
            });
        }
    });
}

function XemTheoKhangNguyen() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();
    $("#btnThemNhanh").hide();
    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/LichSuTiemChung",
        data: { doiTuongId: doiTuongId, viewMode: false },
        type: "GET",
        async: true,
        success: function (response) {
            if (response.Status == 0) {
                jQuery.gritter.add({
                    title: GlobalResources.DT_THONG_BAO,
                    message: response.Message,
                    sticky: false,
                    timeout: '',
                    class_name: 'growl-danger'
                });
            } else {
                $("#lichSuTiem").empty();
                $("#lichSuTiem").html(response);
            }
        },
        error: function (error) {
            jQuery.gritter.add({
                title: GlobalResources.DT_THONG_BAO,
                message: error.message,
                sticky: false,
                timeout: '',
                class_name: 'growl-danger'
            });
        },
        complete: function () {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}

function RefreshLichSuTiemChung() {
    if (ViewMode) {
        XemTheoVacxin();
    } else {
        XemTheoKhangNguyen();
    }
}

function InitTreEmForm() {
    $("#divTreEmArea").show();
    $("#divPhuNuArea").hide();
    $("#divInformationAdvanced").hide();

    $("#slNguoiChamSoc_Detail").select2({ width: "100%", minimumResultsForSearch: -1 });
}

function InitPhuNuForm() {
    $("#divTreEmArea").hide();
    $("#divPhuNuArea").show();
    $("#divInformationAdvanced").show();
    $("#slGioiTinhDetail").select2("val", 1);
    $("#slGioiTinhDetail").prop("disabled", true);
}

function InitDoiTuongKhacForm() {
    $("#divTreEmArea").hide();
    $("#divPhuNuArea").hide();
}

function ShowThemMuiTiemDialog() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();

    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/ThemMuiTiem",
        type: "GET",
        async: true,
        data: { doiTuongId: doiTuongId },
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            $('body').append($(response));
            $('body').on('hidden.bs.modal', '#dgThemMuiTiem', function () {
                $("#dgThemMuiTiem").remove();
            });

            $("#dgThemMuiTiem").modal('show');
            $("#dgThemMuiTiem").draggable();
        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}

function ShowDialogChiTietMuiTiem(lichSuTiemId) {

    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/ChiTietMuiTiem",
        type: "GET",
        data: { lichSuTiemId: lichSuTiemId },
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            $('body').append($(response));
            $('body').on('hidden.bs.modal', '#mdChiTietMuiTiem', function () {
                $("#mdChiTietMuiTiem").remove();
            });

            $("#mdChiTietMuiTiem").modal('show');
            $("#mdChiTietMuiTiem").draggable();
        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}

function ShowDialogCapNhatPhanUng(lichSuTiemId) {
    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/ShowFormCapNhatPhanUng",
        type: "GET",
        data: { lichSuTiemId: lichSuTiemId },
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            $('body').append($(response));
            $('body').on('hidden.bs.modal', '#mdCapNhatPhanUng', function () {
                $("#mdCapNhatPhanUng").remove();
            });

            $("#mdCapNhatPhanUng").modal('show');
            $("#mdCapNhatPhanUng").draggable();

            var loaiPhanUng = $("#LOAI_PHAN_UNG").val();
            var theoDoiSauTiemId = $("#HF_THEO_DOI_SAU_TIEM_ID_CNPU").val();
            if (loaiPhanUng == 2) {
                ShowPhanUngThongThuong(theoDoiSauTiemId);
            } else if (loaiPhanUng == 3) {
                ShowPhanUngNang(theoDoiSauTiemId);
            } else if (loaiPhanUng == 4) {
                ShowPhanUngKhac(theoDoiSauTiemId);
            }
        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}

function ShowDialogSuaMuiTiem(lichSuTiemId) {

    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/SuaMuiTiem",
        type: "GET",
        data: { lichSuTiemId: lichSuTiemId },
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            $('body').append($(response));
            $('body').on('hidden.bs.modal', '#dgSuaMuiTiem', function () {
                $("#dgSuaMuiTiem").remove();
            });

            $("#dgSuaMuiTiem").modal('show');
            $("#dgSuaMuiTiem").draggable();
        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}

function ShowDialogXoaMuiTiem(lichSuTiemId) {

    // hien thi confirm xóa
    bootbox.confirm({
        size: 'small',
        message: GlobalResources.DT_BAN_CO_CHAC_CHAN_MUON_XOA_MUI_TIEM_NAY_KHON + '?',
        title: GlobalResources.DT_XAC_NHAN_XOA_MUI_TIEM,
        callback: function (result) {
            if (result != null) {
                if (result) {
                    AttemptXoaMuiTiem(lichSuTiemId);
                }
            }
        },
        buttons: {
            confirm: {
                label: GlobalResources.DT_DONG_Y,
                className: 'btn-danger'
            },
            cancel: {
                label: GlobalResources.DT_DONG,
                className: 'btn-default'
            }
        }
    });
}

function AttemptXoaMuiTiem(lichSuTiemId) {
    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/XoaMuiTiem",
        type: "POST",
        data: { lichSuTiemId: lichSuTiemId },

        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            if (response.Status == 1) {

                jQuery.gritter.add({
                    title: GlobalResources.DT_THONG_BAO,
                    text: response.Message,
                    sticky: false,
                    timeout: '',
                    class_name: 'growl-success'
                });

                RefreshLichSuTiemChung();
            } else if (response.Status == 0) {

                jQuery.gritter.add({
                    title: GlobalResources.DT_THONG_BAO,
                    text: response.Message,
                    sticky: false,
                    timeout: '',
                    class_name: 'growl-danger'
                });
            }

        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));

            jQuery.gritter.add({
                title: GlobalResources.DT_THONG_BAO,
                text: GlobalResources.DT_PHIEN_DANG_NHAP_CUA_BAN_DA_HET_HAN_KHONG_CO_MANG,
                sticky: false,
                timeout: '',
                class_name: 'growl-danger'
            });
        }
    });
}

function DateSorter(date1, date2) {
    if (CommonValidation.ParseDateStringVi(String(date1).trim()) < CommonValidation.ParseDateStringVi(String(date2).trim())) {
        return 1;
    }

    if (CommonValidation.ParseDateStringVi(String(date1).trim()) > CommonValidation.ParseDateStringVi(String(date2).trim())) {
        return -1;
    }
    return 0;
}

// Thêm số đo dinh dưỡng
//function ThemSoDoDinhDuong(type, rowIndex) {
//    if ($("#tblDinhDuong > tbody > tr.newRowDinhDuong").length > 0) {
//        return;
//    }

//    $("#btnThemSoDoDinhDuong").prop("disabled", true);

//    var rowHtml = '<tr class="newRowDinhDuong"> \
//                            <td style="display:none;"><input name="SoDoCoTheID" id="hfSoDoID" type="hidden" value="" /></td>\
//                            <td class="tdNgayDo" style="text-align: center; vertical-align: middle;"> \
//                                    <input name="NgayDo" type="text" class="dpNgayDo form-control date-picker input-sm" onchange="LoadThangTuoi(this)" placeholder="ngày/tháng/năm" value="" > \
//                            </td> \
//                            <td class="tdThangTuoi" style="text-align: center; vertical-align: middle;"></td> \
//                            <td class="tdCanNang" style="text-align: center; vertical-align: middle;"> \
//                                <input onkeydown="OnlyNumber(event)"  name="canNang" maxlength="7" type="text" class="canNang form-control input-sm" value="" style="text-align: right;" > \
//                            </td> \
//                            <td class="tdChieuCao" style="text-align: center; vertical-align: middle;"> \
//                                <input onkeydown="OnlyNumber(event)" name="chieuCao" maxlength="7" type="text" class="chieuCao form-control input-sm" value="" style="text-align: right;" > \
//                            </td> \
//                            <td class="tdUongSuaMe" style="text-align: center; vertical-align: middle;">\
//                                <input type="checkbox" id="cbUongSuaMe">\
//                            </td>\
//                            <td class="tdUongVitaminA" style="text-align: center; vertical-align: middle;">\
//                                <input type="checkbox" id="cbUongVitaminA">\
//                            </td>\
//                            <td class="tdTrangThai" style="text-align: center; vertical-align: middle;"></td> \
//                            <td class="tdThaoTac" style="text-align: center; vertical-align: middle;"> \
//                                <a href="#" class="fa-hover mr10" title="Lưu thông tin" data-toggle="tooltip" onclick="SaveThongTinhDinhDuong()" > \
//                                    <i class="fa fa-save color-primary"></i> \
//                                </a> \
//                                <a href="#" class="fa-hover" title="Hủy bỏ" data-toggle="tooltip" onclick="HuyLuuThongTinDinhDuong(this)"> \
//                                    <i class="fa fa-reply color-primary"></i> \
//                                </a> \
//                            </td> \
//                      </tr>';

//    var rowObj = $(rowHtml);
//    var colNgayDo = rowObj.find("input[name='NgayDo']").first();

//    jQuery.datetimepicker.setLocale(GlobalResources.DATE_TIME_PICKER_LANGUAGE);   
//    $(colNgayDo).datetimepicker({
//        format: 'd/m/Y',
//        step: 1,
//        i18n: {
//            vi: {
//                months: ["Tháng Một", "Tháng Hai", "Tháng Ba", "Tháng Tư",
//              "Tháng Năm", "Tháng Sáu", "Tháng Bảy", "Tháng Tám", "Tháng Chín",
//              "Tháng Mười", "Tháng Mười Một", "Tháng Mười Hai"],
//                dayOfWeek: ["CN", "T2", "T3", "T4", "T5", "T6", "T7", ]
//            }
//        },
//        mask: '99/99/9999',
//        timepicker: false
//    });
//    $(colNgayDo).mask("99/99/9999");

//    if ($("#tblDinhDuong > tbody > tr.rowEmptyThongTinDinhDuong").length > 0) {
//        $("#tblDinhDuong > tbody > tr.rowEmptyThongTinDinhDuong").remove()
//        $("#tblDinhDuong > tbody").prepend(rowObj);
//        $("#tblDinhDuong > tbody > tr.newRowDinhDuong > td.tdThangTuoi").html("Sơ sinh");
//    } else {
//        if (type == "EDIT") {
//            $('#tblDinhDuong > tbody > tr').eq(rowIndex).after(rowObj);
//        } else if (type = "INSERT") {
//            $("#tblDinhDuong > tbody").prepend(rowObj);
//        }

//        // Disabled các dòng không thao tác
//        rowObj.closest('tr').siblings().css({
//            "pointer-events": "none",
//            "opacity": "0.5"
//        });
//    }
//}

function HuyLuuThongTinDinhDuong(obj) {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();

    LoadThongTinDinhDuong(doiTuongId, PageNumberDD);
}

function ShowModalBieuDoDinhDuongCanNang() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();

    $.ajax({
        url: '/TiemChung/DoiTuong/ShowModalBieuDoDinhDuongCanNang',
        type: "POST",
        data: { doiTuongId: doiTuongId },
        headers: layGiaTriToken(),
        cache: false,
        success: function (response) {

            $("#ModalBieuDoDinhDuongContainer").empty();
            $("#ModalBieuDoDinhDuongContainer").html(response);

            var modalIntance = $('#ModalBieuDoDinhDuongContainer').find('.modal').modal('show');
            modalIntance.on('shown.bs.modal', function () {
                // Recreate the chart now and it will be correct
                $('#chart_bieu_do_dinh_duong_can_nang').highcharts().reflow();
            });

            modalIntance.on('hidden.bs.modal', function () {
                modalIntance.parent().empty();
            });


        },
        error: function (e) {
            CommonJS.showDangerMessage(GlobalResources.DT_PHIEN_DANG_NHAP_CUA_BAN_DA_HET_HAN_KHONG_CO_MANG);
        },
        complete: function (xhr, status) {
            Common.UI.UnBlockElement("#detailPanel");
        }
    });
}

function ShowModalBieuDoDinhDuongChieuCao() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();

    $.ajax({
        url: '/TiemChung/DoiTuong/ShowModalBieuDoDinhDuongChieuCao',
        type: "POST",
        data: { doiTuongId: doiTuongId },
        headers: layGiaTriToken(),
        cache: false,
        success: function (response) {
            $("#ModalBieuDoDinhDuongContainer").empty();
            $("#ModalBieuDoDinhDuongContainer").html(response);

            var modalIntance = $('#ModalBieuDoDinhDuongContainer').find('.modal').modal('show');
            modalIntance.on('shown.bs.modal', function () {
                // Recreate the chart now and it will be correct
                $('#chart_bieu_do_dinh_duong_chieu_cao').highcharts().reflow();
            });

            modalIntance.on('hidden.bs.modal', function () {
                modalIntance.parent().empty();
            });
        },
        error: function (e) {
            CommonJS.showDangerMessage(GlobalResources.DT_PHIEN_DANG_NHAP_CUA_BAN_DA_HET_HAN_KHONG_CO_MANG);
        },
        complete: function (xhr, status) {
            Common.UI.UnBlockElement("#detailPanel");
        }
    });
}

///// Hiển thị modal đăng ký dịch vụ /////
function ShowModalDangKyDichVu() {
    var doiTuongId = TiemChung$DoiTuong$Detail$DoiTuongModel.DOI_TUONG_2_ID;
    Common.UI.BlockElement("#thongTinDichVu");
    $.ajax({
        url: '/TiemChung/DoiTuong/ShowModalDangKyDichVu',
        type: 'GET',
        data: { DoiTuongId: doiTuongId },
        async: true,
        success: function (response) {
            $("#mdlDangKyDichVuContainer").append(response);
            var modalInstance = $("#mdlDangKyDichVuContainer").find(".modal").modal("show");
            modalInstance.on("hidden.bs.modal", function (e) {
                $("#mdlDangKyDichVuContainer").empty();
                // cap nhat danh sach dich vu dang ky
                CapNhatDanhSachDichVuDangKy(doiTuongId);
            });
        },
        error: function () {

        },
        complete: function () {
            Common.UI.UnBlockElement("#thongTinDichVu");
        }
    });
}

function InPhieuThu(dich_vu_dang_ky_id) {
    var baseUrl = '@Url.Action("InPhieuThu", "DoiTuong")'
        + '?DICH_VU_DANG_KY_ID=' + Number(dich_vu_dang_ky_id);
    bootbox.dialog({
        title: GlobalResources.DT_IN_PHIEU_THU,
        message: '<div class="row mt5"><iframe id="frmPrintPreview" src="' + baseUrl + '" height="500" width="100%" border="0"></iframe></div>',
        size: 'large'
    });
}

///// Hiển thị modal lịch sử các tin nhắn đã gửi cho đối tượng /////
function ShowModalLichSuTinNhan() {
    var doiTuongId = TiemChung$DoiTuong$Detail$DoiTuongModel.DOI_TUONG_2_ID;
    Common.UI.BlockElement("#thongTinDichVu");
    $.ajax({
        url: '/TiemChung/DoiTuong/ShowModalLichSuTinNhan',
        type: 'GET',
        data: { DoiTuongId: doiTuongId },
        async: true,
        success: function (response) {
            $("#mdlDangKyDichVuContainer").append(response);
            var modalInstance = $("#mdlDangKyDichVuContainer").find(".modal").modal("show");
            modalInstance.on("hidden.bs.modal", function (e) {
                $("#mdlDangKyDichVuContainer").empty();
            });
        },
        error: function () {

        },
        complete: function () {
            Common.UI.UnBlockElement("#thongTinDichVu");
        }
    });
}


function CapNhatDanhSachDichVuDangKy(doiTuongId) {
    Common.UI.BlockElement("#divDanhSachDichVu");

    $.ajax({
        url: '/TiemChung/DoiTuong/DichVuDangKy',
        type: 'GET',
        data: { DoiTuongId: doiTuongId },
        async: true,
        success: function (response) {
            $("#divDanhSachDichVu").empty();
            $("#divDanhSachDichVu").append(response);
        },
        error: function (e) {
            jQuery.gritter.add({
                text: GlobalResources.DT_KHONG_THE_LAY_DANH_SACH_DICH_VU_DANG_KY_THU_LAI,
                class_name: 'growl-danger',
                timeout: 2000,
                sticky: false
            });
        },
        complete: function () {
            Common.UI.UnBlockElement("#divDanhSachDichVu");
        }
    });
}

function ShowDialogCapNhatSDTDangKyDichVu(dichVuDangKyId) {
    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: "/TiemChung/DoiTuong/ShowFormCapNhatSDTDangKyDichVu",
        type: "GET",
        data: { DichVuDangKyId: dichVuDangKyId },
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            $('body').append($(response));
            $('body').on('hidden.bs.modal', '#dgCapNhatSDT', function () {
                $("#dgCapNhatSDT").remove();
            });

            $("#dgCapNhatSDT").modal('show');
            $("#dgCapNhatSDT").draggable();
        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}

function ShowDialogCapNhatGoiCuoc(dichVuDangKyId) {
    Common.UI.BlockElement($("#detailPanel"));

    $.ajax({
        url: '/TiemChung/DoiTuong/ShowDialogCapNhatGoiCuoc',
        type: "GET",
        data: { dichVuDangKyId: dichVuDangKyId },
        success: function (response) {
            Common.UI.UnBlockElement($("#detailPanel"));

            $('body').append($(response));
            $('body').on('hidden.bs.modal', '#mdCapNhatGoiCuoc', function () {
                $("#mdCapNhatGoiCuoc").remove();
            });

            $("#mdCapNhatGoiCuoc").modal('show');
            $("#mdCapNhatGoiCuoc").draggable();
        },
        error: function (e) {
            Common.UI.UnBlockElement($("#detailPanel"));
        }
    });
}


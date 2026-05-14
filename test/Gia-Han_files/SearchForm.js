$(document).ready(function () {
    InitSearchControls();
}
)

function InitSearchControls() {
    $("#slLoaiDiaChiSearch").select2({ width: "100%", minimumResultsForSearch: -1 });
    $("#slVungMienSearch").select2({ width: "100%", minimumResultsForSearch: -1 });
    $("#slTinhSearch").select2({ width: "100%" });
    $("#slHuyenSearch").select2({ width: "100%" });
    $("#slXaSearch").select2({ width: "100%" });
    $("#slThonApSearch").select2({ width: "100%" });
    $("#slDonViTao").select2({ width: "100%", minimumResultsForSearch: -1 });
    $("#slGioiTinh").select2({ width: '100%', minimumResultsForSearch: -1 });
    $("#slLuaTuoi").select2({ width: '100%', minimumResultsForSearch: -1 });
    $("#slTinhTrangTheoDoi").select2({ width: '100%', minimumResultsForSearch: -1 });
    $("#slTinhTrangMangThai").select2({ width: '100%', minimumResultsForSearch: -1 });
    $("#slVungMienSearch").on("change", function (e) {
        OnVungMienSearchChange(e);
    });

    $("#slTinhSearch").on("change", function () {
        var tinhId = $('#slTinhSearch').select2("val");
        OnTinhSearchChange(tinhId);
    });

    $("#slXaSearch").on("change", function () {
        var xaId = $('#slXaSearch').select2("val");
        OnXaSearchChange(xaId);
    });

    $("#slLuaTuoi").on("change", function (e) {
        OnLuaTuoiSearchChange(e);
    });

    $("#slGioiTinh").on("change", function (e) {
        AnHienTinhTrangMangThai(e);
    });
    jQuery.datetimepicker.setLocale(GlobalResources.DATE_TIME_PICKER_LANGUAGE);
    $("#txtNgayTiemTu").datetimepicker({
        format: 'd/m/Y',
        step: 1,
        timepicker: false,
        defaultDate: getFirstDayOfYear()
    });
    $(".date-picker").datetimepicker({
        format: 'd/m/Y',
        step: 1,
        timepicker: false
    });
    $(".date-picker").mask("99/99/9999");

    $("#btnXuatExcelDsTre").on("click", function () {
        XuatDanhSachDoiTuong();
    });

    $("#btnXuatExcelDsPhuNu").on("click", function () {
        XuatDanhSachDoiTuongPhuNu();
    });

    $("#btnXuatExcelDsThaiPhu").on("click", function () {
        XuatDanhSachDoiTuongThaiPhu();
    });

    $("#btnXuatPDfDsTre").on("click", function () {

    });

    $("#btnXuatExcelSoA21").on("click", function () {
        XuatSoA21();
    });

    $("#btnXuatExcelDSThucHienTiemTaiCSYT").on("click", function () {
        XuatExcelDSThucHienTiemTaiCSYT();
    });

    $("#btnXuatExcelSoA22").on("click", function () {
        //XuatSoA22();
        //$(".modal-backdrop").attr("style", "position:fixed; z-index:1999");
        $('#SelectDateModal').modal('show');
    });

    $("#btnXuatExcelSoA2KhangNguyen").on("click", function () {
        XuatSoA2KhangNguyen();
    });

    $("#btnXuatExcelDSTreTiemChungTrenDiaBan").on("click", function () {
        XuatDanhSachTreTiemChungTrenDiaBan();
    });


    $("#btnXuatExcelSoA23UonVan").on("click", function () {
        XuatSoA23UonVanPhuNu();
    });

    $("#btnXuatExcelDSTiemUonVanChoPhuNu").on("click", function () {
        XuatExcelDSTiemUonVanChoPhuNu();
    });

    $("#btnXuatExcelDangKyDichVu").on("click", function () {
        XuatBieuMauDangKyDichVu();
    });

    // Tự động tìm kiếm khi người dùng nhập mã đối tượng hợp lệ
    $("#txtMaDoiTuongSearch").on("change", function () {
        TuDongTimTheoMaDoiTuong();
    });

    // Tự động tìm kiếm khi người dùng nhấn nút Enter
    $('#frmSearchForm.input').keypress(function (e) {
        if (e.which == 13) {
            $('#frmSearchForm').submit();
            return false;
        }
    });

    KhoiTaoDanhSachTinh();
    KhoiTaoDanhSachXa();
    KhoiTaoDanhSachThonAp();
    KhoiTaoDanhSachDanToc();
    KhoiTaoDanhSachCoSo();
}

function KhoiTaoDanhSachCoSo() {
    var currentLang = MultiLanguage.Cookies.getCookie("LangForMultiLanguage");
    $("#slDonViTao").select2({
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
        "language": currentLang,
        templateResult: function (repo) {
            if (repo.loading) return repo.text;

            var markup = "<div class='select2-result-repository clearfix'> " +
                "<div class='select2-result-repository__meta'> " +
                "<div class='select2-result-repository__title'> "
                + repo.TEN_CO_SO + "</div>"
                + "<div class='select2-result-repository__description'>"
                + (repo.TEN_XA != null ? repo.TEN_XA + ", " : "")
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
}

function BuildSearchCondition() {
    var vungMienId = $("#slVungMienSearch").select2("val");
    var tinhId = $("#slTinhSearch").select2("val");
    var xaId = $("#slXaSearch").select2("val");
    var thonApId = $("#slThonApSearch").select2("val");
    var donViTao = $("#slDonViTao").select2("val");
    var danTocId = $("#slDanTocSearch").select2("val");
    var ngaySinhTu = $("#txtNgaySinhTu").val();
    var ngaySinhToi = $("#txtNgaySinhToi").val();
    var gioiTinh = $("#slGioiTinh").select2("val");
    var luaTuoi = $("#slLuaTuoi").select2("val");
    var maDoiTuong = $("#txtMaDoiTuongSearch").val();
    var hoTen = $("#txtHoTenSearch").val();
    var tenMe = $("#txtTenMeSearch").val();
    var tenBo = $("#txtTenBoSearch").val();
    var maDinhDanh = $("#txtMDDSearch").val();
    var soDienThoai = $("#txtSoDienThoaiSearch").val();
    var tenNguoiGiamHo = $("#txtTenNguoiGiamHoSearch").val();
    var loaiDiaChi = $("#slLoaiDiaChiSearch").val();
    var tinhTrangTheoDoi = $("#slTinhTrangTheoDoi").select2("val");
    var tinhTrangMangThai = $("#slTinhTrangMangThai").select2("val");
    var tuNgay = $('#txtNgayTiemTu').val();
    var toiNgay = $('#txtNgayTiemToi').val();
    return {
        LoaiDiaChi: loaiDiaChi,
        VungMienId: vungMienId,
        TinhId: tinhId,
        XaId: xaId,
        DonViTao: donViTao,
        ThonApId: thonApId,
        DanTocId: danTocId,
        NgaySinhTu: ngaySinhTu,
        NgaySinhToi: ngaySinhToi,
        MaDoiTuong: maDoiTuong,
        TenDoiTuong: hoTen,
        TenBo: tenBo,
        TenMe: tenMe,
        TenNguoiGiamHo: tenNguoiGiamHo,
        MaDinhDanh: maDinhDanh,
        SoDienThoai: soDienThoai,
        LuaTuoi: luaTuoi,
        GioiTinh: gioiTinh,
        TinhTrangTheoDoi: tinhTrangTheoDoi,
        TinhTrangMangThai: tinhTrangMangThai,
        TuNgay: tuNgay,
        ToiNgay: toiNgay
    }
}

function OnCapNhatKhangNguyen() {
    var xaId = $("#slXaSearch").select2("val");
    if (!isNaN(xaId)) {
        Common.UI.BlockElement($("body"));

        $.ajax({
            type: "POST",
            url: "/TiemChung/DoiTuong/CapNhatKhangNguyen",
            data: { xaId: xaId },

            success: function (response) {
                Common.UI.UnBlockElement($("body"));
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: GlobalResources.DT_CAP_NHAT_KHANG_NGUYEN_THANH_CONG,
                    sticky: false,
                    timeout: '',
                    class_name: 'growl-success'
                })
            },
            error: function (e) {
                Common.UI.UnBlockElement($("body"));
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: GlobalResources.DT_LOI_CAP_NHAT_KHANG_NGUYEN,
                    sticky: false,
                    timeout: '',
                    class_name: 'growl-danger'
                });
            }
        });
    }
}
function safeRedirect(fullUrl) {
    try {
        var absolute = fullUrl.startsWith('/') ?
            window.location.origin + fullUrl :
            fullUrl;
        var u = new URL(absolute);
        if (u.origin === window.location.origin &&
            u.pathname.indexOf('/TiemChung/') === 0) {
            window.location.href = u.toString();
            return true;
        }
    } catch (e) { }
    console.warn("Redirect bị chặn do không hợp lệ:", fullUrl);
    return false;
}
function XuatDanhSachDoiTuong() {
    if (ValidateSearchFormXuatDanhSach(1, false)) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatDanhSachDoiTuong' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatDanhSachDoiTuong";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}

function XuatDanhSachDoiTuongPhuNu() {
    if (ValidateSearchFormXuatDanhSach(1, false)) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatDanhSachDoiTuongPhuNu' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatDanhSachDoiTuongPhuNu";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}

function XuatDanhSachDoiTuongThaiPhu() {
    if (ValidateSearchFormXuatDanhSach(1, false)) {
        var lang = $("#LANG").val() ?? "en";

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        CommonJS.downloadFile("GET", '/TiemChung/DoiTuong/XuatDanhSachDoiTuongThaiPhu' + '?' + urlParams)
            .catch(e => {
                CommonJS.showCommonDangerMessageWithLang(lang);
            });
    }
}

/**
 * Xuất danh sách đối tượng theo biểu mẫu sổ A21
 */
function XuatSoA21() {
    if (ValidateSearchFormXuatDanhSach(0, false)) {
        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatSoA21' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatSoA21";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}
function XuatExcelDSThucHienTiemTaiCSYT() {
    if (ValidateSearchForm(1)) {
        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreDSThucHienTiemTaiCSYT' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/DSThucHienTiemTaiCSYT";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}
function ValidateXuatSoA22() {
    var today = new Date();
    var tuNgay = $('#txtNgayTiemTu').val();
    var toiNgay = $('#txtNgayTiemToi').val();

    if (tuNgay == "" || toiNgay == "") {
        CommonJS.showDangerMessage("Bạn phải nhập thực hiện từ ngày và tới ngày.");
        return false;
    }
    if (tuNgay != "" && !CommonJS.ValidateDate(tuNgay)) {
        CommonJS.showDangerMessage("Thực hiện từ ngày không đúng định dạng.");
        return false;
    }
    if (toiNgay != "" && !CommonJS.ValidateDate(toiNgay)) {
        CommonJS.showDangerMessage("Tới ngày không đúng định dạng.");
        return false;
    }
    if (CommonJS.ckFormatDate(tuNgay) > CommonJS.ckFormatDate(toiNgay)) {
        CommonJS.showDangerMessage("Thực hiện từ ngày không được lớn hơn tới ngày.");
        return false;
    }
    if (processDate(tuNgay) > today) {
        CommonJS.showDangerMessage("Thực hiện từ ngày và tới ngày không được lớn hơn ngày hiện tại.");
        return false;
    }
    if (processDate(toiNgay) > today) {
        CommonJS.showDangerMessage("Thực hiện từ ngày và tới ngày không được lớn hơn ngày hiện tại.");
        return false;
    }
    if (processDate(tuNgay).getFullYear() != processDate(toiNgay).getFullYear()) {
        CommonJS.showDangerMessage("Thực hiện từ ngày và tới ngày phải nằm trong cùng một năm.");
        return false;
    }
    return true;
}

function XuatSoA22() {
    if (ValidateSearchForm(1, false) == true && ValidateXuatSoA22() == true) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatSoA22' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatSoA22";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            },
            complete: function () {
                $('#SelectDateModal').modal('hide');
            }
        });
    }
}

function XuatSoA2KhangNguyen() {
    if (ValidateSearchForm(0, false)) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatSoA2KhangNguyen' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatSoA2KhangNguyen";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}

function XuatDanhSachTreTiemChungTrenDiaBan() {
    if (ValidateSearchForm(0, false)) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatDanhSachTreTiemChungTrenDiaBan' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatDanhSachTreTiemChungTrenDiaBan";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}

function XuatSoA23UonVanPhuNu() {
    if (ValidateSearchForm(0, false)) {
        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');
        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatSoA23UonVanPhuNu' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatSoA23UonVanPhuNu";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}

function XuatExcelDSTiemUonVanChoPhuNu() {
    if (ValidateSearchForm(0)) {
        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');
        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatDSTiemUonVanChoPhuNu' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatDSTiemUonVanChoPhuNu";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}


function XuatBieuMauDangKyDichVu() {
    if (ValidateSearchForm(1, false)) {
        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');
        $.ajax({
            url: '/TiemChung/DoiTuong/PreXuatBieuMauDangKyDichVu' + '?' + urlParams,
            async: true,
            type: "GET",
            success: function (result) {
                if (result.Status == 0) {
                    jQuery.gritter.add({
                        text: result.Message,
                        class_name: 'growl-danger',
                        timeout: 2000,
                        sticky: false
                    });
                } else if (result.Status == 1) {
                    var basePath = "/TiemChung/DoiTuong/XuatBieuMauDangKyDichVu";
                    var finalUrl = basePath + "?" + urlParams;
                    safeRedirect(finalUrl);
                }
            }
        });
    }
}



function OnVungMienSearchChange(e) {
    var vungMienId = $("#slVungMienSearch").select2("val");
    if (isNaN(vungMienId)) {
        vungMienId = -1;
    }

    ImCache.GetDsTinhByVungMienFromCache(vungMienId, function (dsTinh) {
        CapNhatDsTinhSearch(dsTinh);
    });
}

//===== BEGIN: Xu ly thay doi danh muc tinh, huyen, xa =====//
function KhoiTaoDanhSachTinh() {
    var vungMienId = TiemChung$DoiTuong$SearchFormModel.VungMienId;
    if (vungMienId == null) {
        vungMienId = -1;
    }

    ImCache.GetDsTinhByVungMienFromCache(vungMienId, function (dsTinh) {
        CapNhatDsTinhSearch(dsTinh);

        // thiết lập tỉnh được chọn mặc định
        $("#slTinhSearch").val(TiemChung$DoiTuong$SearchFormModel.TinhId);
        $("#slTinhSearch").select2({});
    });
}

function KhoiTaoDanhSachXa() {
    var tinhId = TiemChung$DoiTuong$SearchFormModel.TinhId;
    if (!isNaN(tinhId)) {
        ImCache.GetDsXaByTinhFromCache(tinhId, function (dsXa) {
            CapNhatDsXaSearch(dsXa);

            // thiết lập xã được chọn mặc định
            $("#slXaSearch").val(TiemChung$DoiTuong$SearchFormModel.XaId);
            $("#slXaSearch").select2({});
        });
    } else {
        $("#slXaSearch").empty();
        $("#slThonApSearch").empty();
    }

    $("#slXaSearch").select2();
    $("#slThonApSearch").select2();
}

function KhoiTaoDanhSachThonAp() {
    var slThonAp = $("#slThonApSearch");

    slThonAp.empty();
    slThonAp.append($("<option selected='selected'>-" + GlobalResources.DT_THON_AP + "-</option>"));

    if (TiemChung$DoiTuong$SearchFormModel.DsThonAp != null) {
        $.each(TiemChung$DoiTuong$SearchFormModel.DsThonAp, function (index, item) {
            var option = $("<option />");
            option.val(item.THON_AP_ID);
            option.text(item.TEN_THON_AP);
            slThonAp.append(option);
        });
    }

    $("#slThonApSearch").select2();
}

function CapNhatDsTinhSearch(dsTinh) {
    var slTinh = $("#slTinhSearch");
    var slXa = $("#slXaSearch");
    var slThonAp = $("#slThonApSearch");

    slTinh.empty();
    slTinh.append($("<option selected='selected'>-" + GlobalResources.DT_TINH_THANH_PHO + "-</option>"));

    slXa.empty();
    slThonAp.empty();

    $.each(dsTinh, function (index, item) {
        var option = $("<option />");
        option.val(item.TINH_ID);
        option.text(item.TENTINH);
        slTinh.append(option);
    });

    $("#slTinhSearch").select2();
    $("#slXaSearch").select2();
    $("#slThonApSearch").select2();
}

function OnTinhSearchChange(tinhId) {

    if (!isNaN(tinhId)) {
        ImCache.GetDsXaByTinhFromCache(tinhId, function (dsXa) {
            CapNhatDsXaSearch(dsXa);
        });
    } else {
        $("#slXaSearch").empty();
        $("#slThonApSearch").empty();
    }

    $("#slXaSearch").select2();
    $("#slThonApSearch").select2();
}

function CapNhatDsXaSearch(dsXa) {
    var slXa = $("#slXaSearch");
    var slThonAp = $("#slThonApSearch");

    slXa.empty();
    slXa.append($("<option selected='selected'>-" + GlobalResources.DT_XA_PHUONG + "-</option>"));

    slThonAp.empty();

    $.each(dsXa, function (index, item) {
        var option = $("<option />");
        option.val(item.XA_ID);
        option.text(item.TEN_XA);
        slXa.append(option);
    });
}

function OnXaSearchChange(xaId) {
    if (!isNaN(xaId)) {
        $.ajax({
            url: '/DonViHanhChinh/DsThonAp',
            type: "GET",
            data: { xaId: xaId },
            success: function (response) {
                var slThonAp = $("#slThonApSearch");

                slThonAp.empty();
                slThonAp.append($("<option selected='selected'>-" + GlobalResources.DT_THON_AP + "-</option>"));

                $.each(response, function (index, item) {
                    var option = $("<option />");
                    option.val(item.THON_AP_ID);
                    option.text(item.TEN_THON_AP);
                    slThonAp.append(option);
                });
            }
        });
    } else {
        $("#slThonApSearch").empty();
    }

    $("#slThonApSearch").select2();
}
//===== END: Xu ly thay doi danh muc tinh, xa =====//

function KhoiTaoDanhSachDanToc() {
    ImCache.GetDsDanTocFromCache(function (dsDanToc) {
        $.fn.select2.amd.require(['select2/compat/matcher'], function (oldMatcher) {
            $("#slDanTocSearch").select2({
                data:
                    $.map(dsDanToc, function (item) {
                        return {
                            id: item.DAN_TOC_ID,
                            text: item.TEN_DAN_TOC,
                            TEN_DAN_TOC: item.TEN_DAN_TOC,
                            TEN_GOI_KHAC: item.TEN_GOI_KHAC
                        };
                    }),

                placeholder: '-' + GlobalResources.DT_DAN_TOC + '-',
                matcher: oldMatcher(function (term, text, option) {
                    if (option.TEN_DAN_TOC.toUpperCase().indexOf(term.toUpperCase()) > -1 || (option.TEN_GOI_KHAC != null && option.TEN_GOI_KHAC.toUpperCase().indexOf(term.toUpperCase()) > -1)) {
                        return true;
                    }

                    return false;
                }),

                templateResult: function (repo) {
                    if (repo.loading) return repo.text;

                    var markup = "<div class='select2-result-repository clearfix'>\
                                <div class='select2-result-repository__meta'>\
                                    <div class='select2-result-repository__title'>"
                        + repo.TEN_DAN_TOC + "</div>"
                        + "<div class='select2-result-repository__description'>"
                        + (repo.TEN_GOI_KHAC != null ? repo.TEN_GOI_KHAC + ", " : "")
                        + "</div><div><div>";

                    return markup;
                },
                templateSelection: function (repo) {
                    return repo.text || repo.TEN_DAN_TOC;
                },
                escapeMarkup: function (markup) { return markup; },
                allowClear: true
            });

            $("#slDanTocSearch").val(null).trigger("change");
        });
    });
}

function OnLuaTuoiSearchChange(e) {
    var luatuoi = $("#slLuaTuoi").select2("val");
    if (isNaN(luatuoi) || luatuoi == -1) { // khong chon lua tuoi
        // enable phan gioi tinh
        $("#slGioiTinh").val(-1).trigger("change");
        $("#slGioiTinh").prop("disabled", false);

        $("#txtTenMeSearch").parent().hide();
        $("#txtTenBoSearch").parent().hide();
        $("#txtTenNguoiGiamHoSearch").parent().parent().hide();

        $("#txtTenMeSearch").val('');
        $("#txtTenBoSearch").val('');
        $("#txtTenNguoiGiamHoSearch").val('');

        $("#txtMDDSearch").parent().show();
        $("#txtSoDienThoaiSearch").parent().show();

    } else if (luatuoi == 1) { // tre em
        $("#slGioiTinh").val(-1).trigger("change");
        $("#slGioiTinh").prop("disabled", false);

        $("#txtTenMeSearch").parent().show();
        $("#txtTenBoSearch").parent().show();
        $("#txtTenNguoiGiamHoSearch").parent().parent().show();

        $("#txtMDDSearch").parent().hide();
        $("#txtSoDienThoaiSearch").parent().hide();

        $("#txtMDDSearch").val('');
        $("#txtSoDienThoaiSearch").val('');


    } else if (luatuoi == 2) { // phu nu
        $("#slGioiTinh").val(1).trigger("change");
        $("#slGioiTinh").prop("disabled", true);

        $("#txtTenMeSearch").parent().hide();
        $("#txtTenBoSearch").parent().hide();
        $("#txtTenNguoiGiamHoSearch").parent().parent().hide();

        $("#txtMDDSearch").parent().show();
        $("#txtSoDienThoaiSearch").parent().show();

        $("#txtTenMeSearch").val('');
        $("#txtTenBoSearch").val('');
        $("#txtTenNguoiGiamHoSearch").val('');
    }

    AnHienTinhTrangMangThai(e);
}

function AnHienTinhTrangMangThai(e) {
    $("#slTinhTrangMangThai").val(-1).trigger("change");
    $("#slTinhTrangMangThai").parent().hide();
    var luatuoi = $("#slLuaTuoi").select2("val");
    var gioitinh = $("#slGioiTinh").select2("val");
    if ((isNaN(luatuoi) || luatuoi == -1) && gioitinh == 1) { // khong chon lua tuoi va gioi tinh nu
        $("#slTinhTrangMangThai").parent().show();
    }
    else
        if (luatuoi == 2 && gioitinh == 1) { // phu nu
            $("#slTinhTrangMangThai").parent().show();
        }
}

function TuDongTimTheoMaDoiTuong() {
    if ($("#txtMaDoiTuongSearch").val().length >= 15) {
        $("#frmSearchForm").submit();
    }
}

function OnBeginSearch() {

    typeSearch = TypeAdvancedSearch;

    if (ValidateSearchForm(1)) {
        var searchResult = $("#searchResult");
        // xoa du lieu truoc khi thuc hien tim kiem
        searchResult.empty();
        // show loading
        var loading = $("<div id='loading' style='margin:10px; text-align:center;'><span class='fa fa-spinner fa-spin fa-2x fa-fw'></span></div>")
        searchResult.append(loading);
        return true;
    } else {
        return false;
    }
}

function ValidateSearchForm(limitSearchDateCondition, requireNameAndGender = true) {
    // clear all error message
    ClearAllErrorMessages();
    var isValid = true;
    isValid = ValidateNgayThang(limitSearchDateCondition);

    isValid = validateThongTinHC(requireNameAndGender) && isValid;
    return isValid;
}
function ValidateSearchFormXuatDanhSach(limitSearchDateCondition, requireNameAndGender = true) {
    // clear all error message
    ClearAllErrorMessages();
    var isValid = true;
    isValid = ValidateNgayThang(limitSearchDateCondition);
    return isValid;
}
function OnSearchFailed() {
    // show notification

}

function OnInMaVachClick() {
    if (ValidateSearchForm(1)) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');

        window.open('/TiemChung/DoiTuong/InMaVach?' + urlParams, GlobalResources.TT_DANH_SACH_DOI_TUONG, 'width=800,height=1280,scrollbars=yes');
    }
}

function InMaVachHangLoat() {
    if (ValidateSearchForm(1)) {

        var searchCondition = BuildSearchCondition();
        var urlParams = Object.keys(searchCondition).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(searchCondition[key]);
        }).join('&');
        var baseUrl = '/TiemChung/DoiTuong/InMaVachHangLoat?LoaiDiaChi=' + searchCondition.LoaiDiaChi
            + '&VungMienId=' + searchCondition.VungMienId
            + '&TinhId=' + searchCondition.TinhId
            + '&XaId=' + searchCondition.XaId
            + '&DonViTao=' + searchCondition.DonViTao
            + '&ThonApId=' + searchCondition.ThonApId
            + '&DanTocId=' + searchCondition.DanTocId
            + '&NgaySinhTu=' + searchCondition.NgaySinhTu
            + '&NgaySinhToi=' + searchCondition.NgaySinhToi
            + '&MaDoiTuong=' + searchCondition.MaDoiTuong
            + '&TenDoiTuong=' + searchCondition.TenDoiTuong
            + '&TenBo=' + searchCondition.TenBo
            + '&TenMe=' + searchCondition.TenMe
            + '&TenNguoiGiamHo=' + searchCondition.TenNguoiGiamHo
            + '&MaDinhDanh=' + searchCondition.MaDinhDanh
            + '&SoDienThoai=' + searchCondition.SoDienThoai
            + '&LuaTuoi=' + searchCondition.LuaTuoi
            + '&GioiTinh=' + searchCondition.GioiTinh
            + '&TinhTrangMangThai=' + searchCondition.TinhTrangMangThai;
        bootbox.dialog({
            title: GlobalResources.DT_IN_MA_VACH_HANG_LOAT,
            message: '<div class="row mt5"><iframe id="frmPrintPreview" src="' + baseUrl + '" height="500" width="100%" border="0"></iframe></div>',
            size: 'large'
        });
    }
}

function ClearAllErrorMessages() {
    $("#frmSearchForm").find(".error-message").hide();
}

function ValidateNgayThang(x_) {
    var NgaySinhTu = $("#txtNgaySinhTu").val();
    var NgaySinhToi = $("#txtNgaySinhToi").val();

    if (NgaySinhTu.trim() == '') {
        $("#ngaySinhTuValidate").text("Vui lòng nhập Ngày sinh từ");
        $("#ngaySinhTuValidate").closest(".error-message").show();
        return false;
    }

    if (NgaySinhToi.trim() == '') {
        $("#ngaySinhToiValidate").text("Vui lòng nhập Ngày sinh đến");
        $("#ngaySinhToiValidate").closest(".error-message").show();
        return false;
    }

    if (NgaySinhTu.trim() != '' && !ValidateDate(NgaySinhTu)) {
        $("#ngaySinhTuValidate").text(GlobalResources.DT_Dinhdangngaythangkhonghople);
        $("#ngaySinhTuValidate").closest(".error-message").show();
        return false;
    }

    if (NgaySinhToi.trim() != '' && !ValidateDate(NgaySinhToi)) {
        $("#ngaySinhToiValidate").text(GlobalResources.DT_Dinhdangngaythangkhonghople);
        $("#ngaySinhToiValidate").closest(".error-message").show();
        return false;
    }

    if (NgaySinhTu.trim() != '' && CommonJS.ckFormatDate(NgaySinhTu) > CommonJS.ckFormatDate($("#CurrentSystemDate").val())) {
        $("#ngaySinhTuValidate").text(GlobalResources.DT_Ngaysinhtukhongduoclonhonngayhientai);
        $("#ngaySinhTuValidate").closest(".error-message").show();
        return false;
    }

    if (NgaySinhToi.trim() != '' && CommonJS.ckFormatDate(NgaySinhToi) > CommonJS.ckFormatDate($("#CurrentSystemDate").val())) {
        $("#ngaySinhToiValidate").text(GlobalResources.DT_Ngaysinhtoikhongduoclonhonngayhientai);
        $("#ngaySinhToiValidate").closest(".error-message").show();
        return false;
    }

    if ((NgaySinhTu.trim() != '' && NgaySinhToi.trim() != '') && CommonJS.ckFormatDate(NgaySinhTu) > CommonJS.ckFormatDate(NgaySinhToi)) {
        $("#ngaySinhTuValidate").text(GlobalResources.DT_Ngaysinhtukhongduoclonhonngaysinhtoi);
        $("#ngaySinhTuValidate").closest(".error-message").show();
        return false;
    }

    if (Number(x_) == 0) {
        if (NgaySinhTu.trim() == '' && NgaySinhToi.trim() != '') {
            $("#ngaySinhTuValidate").text(GlobalResources.DT_NHAP_NGAY_SINH_TU);
            $("#ngaySinhTuValidate").closest(".error-message").show();
            return false;
        }

        if (NgaySinhTu.trim() != '' && NgaySinhToi.trim() == '') {
            $("#ngaySinhToiValidate").text(GlobalResources.DT_NHAP_NGAY_SINH_TOI);
            $("#ngaySinhToiValidate").closest(".error-message").show();
            return false;
        }

        if (NgaySinhTu.trim() !== '' && NgaySinhToi.trim() !== '') {
            var fromDate = new Date(NgaySinhTu);
            var toDate = new Date(NgaySinhToi);
            var twoYearsLater = new Date(fromDate);
            twoYearsLater.setFullYear(twoYearsLater.getFullYear() + 2);

            if (toDate > twoYearsLater) {
                $("#ngaySinhTuValidate").text(GlobalResources.DT_NGAY_SINH_TOI_VA_NGAY_SINH_TU_TRONG_KHOANG_HAI_NAM);
                $("#ngaySinhTuValidate").parent().show();
                return false;
            }
        }
    }

    return true;
}

function validateThongTinHC(requireNameAndGender = true) {
    let isValid = true;

    const loaiDiaChi = $("#slLoaiDiaChiSearch").val();
    if (!loaiDiaChi || isNaN(loaiDiaChi) || parseInt(loaiDiaChi) < 0) {
        $("#loaiDiaChiValidate").text("Vui lòng chọn Loại địa chỉ hợp lệ.");
        $("#loaiDiaChiValidate").closest(".error-message").show();
        isValid = false;
    }

    const vungMien = $("#slVungMienSearch").val();
    if (!vungMien || isNaN(vungMien) || parseInt(vungMien) <= 0) {
        $("#vungMienValidate").text("Vui lòng chọn Khu vực hợp lệ.");
        $("#vungMienValidate").closest(".error-message").show();
        isValid = false;
    }

    const tinh = $("#slTinhSearch").val();
    if (!tinh || tinh.trim() === "" || tinh === "-1" || isNaN(tinh) || parseInt(tinh) <= 0) {
        $("#tinhValidate").text("Vui lòng chọn Tỉnh / Thành phố hợp lệ.");
        $("#tinhValidate").closest(".error-message").show();
        isValid = false;
    }

    const xa = $("#slXaSearch").val();
    if (!xa || xa.trim() === "" || xa === "-1" || isNaN(xa) || parseInt(xa) <= 0) {
        $("#xaValidate").text("Vui lòng chọn Xã / Phường hợp lệ.");
        $("#xaValidate").closest(".error-message").show();
        isValid = false;
    }

    if (requireNameAndGender) {
        const hoTen = $("#txtHoTenSearch").val();
        if (!hoTen || hoTen.trim() === "") {
            $("#hoTenValidate").text("Vui lòng nhập Họ và tên.");
            $("#hoTenValidate").closest(".error-message").show();
            isValid = false;
        }

        const gioiTinh = $("#slGioiTinh").val();
        if (
            !gioiTinh ||
            gioiTinh.trim() === "" ||
            gioiTinh === "-1" ||
            isNaN(gioiTinh) ||
            ![0, 1, 2].includes(parseInt(gioiTinh)) // không phải 0 hoặc 1
        ) {
            $("#gioiTinhValidate").text("Vui lòng chọn Giới tính hợp lệ.");
            $("#gioiTinhValidate").closest(".error-message").show();
            isValid = false;
        }
    }

    return isValid;
}

const ConvertToDate = s => {

    let [d, m, y] = s.split(/\D/);

    return new Date(y, m - 1, d);
};

function DateDiff(from, to) {
    return (ConvertToDate(to) - ConvertToDate(from)) / 1000 / 24 / 60 / 60
}

function processDate(date) {
    var parts = date.split("/");
    return new Date(parts[2], parts[1] - 1, parts[0]);
}

function ValidateExportScope() {
    // chi duoc phep xuat du lieu do don vi quan ly
    var loaiDiaChi = $("#slLoaiDiaChiSearch").val();

    if (loaiDiaChi != 2) {
        $("#loaiDiaChiValidate").text(GlobalResources.DT_BAN_CHI_XUAT_DT_THEO_DIA_CHI_DK_TIEM);
        $("#loaiDiaChiValidate").parent().show();
        return false;
    }

    // chi duoc phep xuat du lieu cua 1 xa
    var xaId = $("#slXaSearch").val();
    if (xaId == null || xaId <= 0) {
        $("#xaValidate").text(GlobalResources.DT_BAN_PHAI_CHON_XA);
        $("#xaValidate").parent().show();

        return false;
    } else {
        if (xaId != TiemChung$DoiTuong$SearchFormModel.XaId) {
            $("#xaValidate").text(GlobalResources.DT_BAN_KHONG_DUOC_XUAT_DL_CUA_CS_KHAC);
            $("#xaValidate").parent().show();
            return false;
        }
    }

    return true;
}

function ValidateDate(DateString) {
    var regex = /^(((0[1-9]|[12]\d|3[01])\/(0[13578]|1[02])\/((1[6-9]|[2-9]\d)\d{2}))|((0[1-9]|[12]\d|30)\/(0[13456789]|1[012])\/((1[6-9]|[2-9]\d)\d{2}))|((0[1-9]|1\d|2[0-8])\/02\/((1[6-9]|[2-9]\d)\d{2}))|(29\/02\/((1[6-9]|[2-9]\d)(0[48]|[2468][048]|[13579][26])|((16|[2468][048]|[3579][26])00))))$/;
    if (!(regex.test(DateString))) {
        return false;
    }
    else {
        return true;
    }
}

function getFirstDayOfYear() {
    var today = new Date();
    var dd = '01';
    var mm = '01';
    var yyyy = today.getFullYear();
    return dd + '/' + mm + '/' + yyyy;
}

function DatLaiFormTimKiem() {
    //Loai dia chi
    if ($("#LOAI_DIA_CHI_ID").val() == 2) {
        $("#slLoaiDiaChiSearch").val(2);
    } else {
        $("#slLoaiDiaChiSearch").val(0);
    }
    $("#slLoaiDiaChiSearch").change();

    //Khu vuc
    if ($("#VUNG_MIEN_ID").val() != null) {
        $("#slVungMienSearch").val($("#VUNG_MIEN_ID").val());
    }
    $("#slVungMienSearch").change();

    //Tinh, huyen, xa
    KhoiTaoDanhSachTinh();
    KhoiTaoDanhSachXa();
    KhoiTaoDanhSachThonAp();
    KhoiTaoDanhSachDanToc();
    KhoiTaoDanhSachCoSo();

    //Don vi tao
    if ($("#DON_VI_TAO_ID").val() != null) {
        $("#slDonViTao").val($("#DON_VI_TAO_ID").val());
    } else {
        $("#slDonViTao").val('');
    }
    $("#slDonViTao").change();

    $("#txtNgaySinhTu").val('');
    $("#txtNgaySinhToi").val('');

    $("#slGioiTinh").val(-1);
    $("#slGioiTinh").change();

    $("#slLuaTuoi").val(-1);
    $("#slLuaTuoi").change();

    $("#txtMaDoiTuongSearch").val('');
    $("#txtHoTenSearch").val('');
    $("#txtTenMeSearch").val('');
    $("#txtTenBoSearch").val('');
    $("#txtMDDSearch").val('');
    $("#txtSoDienThoaiSearch").val('');
    $("#txtTenNguoiGiamHoSearch").val('');

    $("#slTinhTrangMangThai").val(-1);
    $("#slTinhTrangMangThai").change();
}

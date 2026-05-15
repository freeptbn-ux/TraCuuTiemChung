var TypeQuickSearch = 1;
var TypeAdvancedSearch = 2;
var typeSearch = 1; // Khoi tao se la tim kiem nang cao
var TypeAddTreEm = 1; //
var TypeAddPhuNu = 2; //
var TypeAddKhac = 3;

var typeAdd = TypeAddTreEm;
var PageNumberDD = 1;
//
//----------- Index ---------------------//
$(".children").find(".active").eq(1).removeClass("active");
$(document).on("click", "#btnQuickSearch", function () {
    $("#hfQuickSearchPageNumber").val(1);
    $("#hfQuickSearchPageSize").val(20);
    QuickSearch();
}).on("click", "#btnThemPhuNu", function () {
    ShowAddPhuNuForm();
}).on("click", "#btnThemTreEm", function () {
    ShowAddTreEmForm();
}).on("click", "#btnAdd", function () {
    ShowAddDoiTuongKhacForm();
}).on("click", "#btnEdit", function () {
    ShowEditForm();
}).on("click", "#btnDelete", function () {
    ConfirmDelete();
}).on("click", "#btnSave", function () {
    Save();
}).on("click", "#btnCancel", function () {
    Cancel();
}).on("click", "#btnAdd2", function () {
    ShowAddDoiTuongKhacForm();
}).on("click", "#btnEdit2", function () {
    ShowEditForm();
}).on("click", "#btnDelete2", function () {
    ConfirmDelete();
}).on("click", "#btnSave2", function () {
    Save();
}).on("click", "#btnCancel2", function () {
    Cancel();
}).on("click", "#btnPrintBarcode", function () {
    var MaDoiTuong = $('#txtMaDoiTuong').val();
    var HoTen = $('#txtHoTen').val();
    PrintBarcode(MaDoiTuong, HoTen);
}).on("click", "#btnAddDTFromFile", function () {
    ShowFormChooseFileAddDoiTuong();
}).on("click", "#btnMigrate", function () {
    MigrateDuLieu();
});

$("#txtKeyword").keypress(function (event) {
    if (event.which == 13) {
        event.preventDefault();
        $("#hfQuickSearchPageNumber").val(1);
        $("#hfQuickSearchPageSize").val(20);
        QuickSearch();
    }
});

function OnShowPage(pageNumber, pageSize) {
    if (typeSearch == TypeQuickSearch) {
        $("#hfQuickSearchPageNumber").val(pageNumber);
        $("#hfQuickSearchPageSize").val(pageSize);
        QuickSearch();

    } else {
        $("#hfSearchFormPageNumber").val(pageNumber);
        $("#hfSearchFormPageSize").val(pageSize);

        // re-submit search form
        var btnAdvancedSearch = $("#btnAdvancedSearch");
        btnAdvancedSearch.submit();
    }
}

function SendXacMinhThongTin(__value) {
    var doiTuongId = __value.trim();
    bootbox.confirm({
        size: 'small',
        message: 'Bạn muốn gửi thông tin sang hệ thống Hồ sơ sức khỏe điện tử?',
        title: GlobalResources.DT_THONG_BAO,
        callback: function (result) {
            if (result != null) {
                if (result) {
                    SendXacMinhThongTinDoiTuong(doiTuongId);
                }
            }
        },
        buttons: {
            confirm: {
                label: GlobalResources.DONG_Y,
                className: 'btn-warning btn-sm'
            },
            cancel: {
                label: GlobalResources.DONG,
                className: 'btn-default btn-sm'
            }
        }
    });
}

function SendXacMinhThongTinDoiTuong(doiTuongId) {
    Common.UI.BlockElement("#detailPanel");
    $.ajax({
        type: "POST",
        url: "/TiemChung/DoiTuong/Covid_SendXacMinhThongTinDoiTuong",
        data: { doiTuongId_str: doiTuongId },
        headers: layGiaTriToken(),
        async: true,
        success: function (response) {
            Common.UI.UnBlockElement("#detailPanel");
            if (response.Status == 0) {
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: response.Message,
                    class_name: 'growl-danger',
                    timeout: '',
                    sticky: false
                });
            } else if (response.Status == 1) {
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: response.Message,
                    class_name: 'growl-success',
                    timeout: '',
                    sticky: false
                });

                // refresh lai danh sach tim kiem
                OnShowPage(1, 20);

                var doiTuongDetail = $("#doiTuongDetailCovid");
                doiTuongDetail.empty();

                var message = $("<div style='padding: 10px'>" + GlobalResources.DT_CHON_MOT_DOI_TUONG_DE_XEM_THONG_TIN + "</div>");
                doiTuongDetail.append(message);
                ShowActionButtonDefault();
            }
        },
        error: function (e) {
            Common.UI.UnBlockElement("#detailPanel");

            jQuery.gritter.add({
                title: GlobalResources.THONG_BAO,
                text: e.message,
                class_name: 'growl-danger',
                timeout: '',
                sticky: false
            });
        }
    });
}

function ShowFormChooseFileAddDoiTuong() {

    $.ajax({
        type: "POST",
        url: '/TiemChung/DoiTuong/ThemDSDoiTuongTuExcel',
        catche: false,
        traditional: true,
        dataType: "html",
        success: function (data) {
            $("#ModalDialogThemDT").html(data)
            $("#ModalDialogThemDT").modal("show");

        }
    });
}

function QuickSearch() {
    typeSearch = TypeQuickSearch;

    var keyword = CommonJS.htmlEncode($("#txtKeyword").val().trim());
    var searchResult = $("#searchResult");
    Common.UI.BlockElement("#searchResult");

    var pageNumber = $("#hfQuickSearchPageNumber").val();
    var pageSize = $("#hfQuickSearchPageSize").val();

    $.ajax({
        url: '/TiemChung/DoiTuong/TimKiemNhanh',
        type: "POST",
        async: true,
        data: { keyword: keyword, pageSize: pageSize, pageNumber: pageNumber },
        success: function (response) {
            searchResult.empty();
            searchResult.html(response);
        },
        complete: function (xhr, status) {
            Common.UI.UnBlockElement("#searchResult");
        }
    });
}

function ShowAddPhuNuForm() {
    typeAdd = TypeAddPhuNu;
    ShowAddForm();
}

function ShowAddTreEmForm() {
    typeAdd = TypeAddTreEm;
    ShowAddForm();
}

function ShowAddDoiTuongKhacForm() {
    typeAdd = TypeAddKhac;
    ShowAddForm();
}

function ShowAddForm() {

    var doiTuongDetail = $("#doiTuongDetail");
    Common.UI.BlockElement("#detailPanel");

    $.ajax({
        url: '/TiemChung/DoiTuong/AddNew',
        type: "GET",
        async: true,
        data: null,
        success: function (response) {
            doiTuongDetail.empty();
            doiTuongDetail.html(response);
            var hfTypeAdd = $("#hfTypeAdd_ThemMoi");
            hfTypeAdd.val(typeAdd);

            InitFormThemMoi(); // invoke function from AddNew.cshtml
            ShowActionButtonForAddNew();
        },
        error: function (e) {
            doiTuongDetail.empty();

            // show error notification
            var errorMessageElement = $("<div />");
            errorMessageElement.attr("style", "margin:10px; text-align:center;");
            errorMessageElement.addClass("has-error");
            errorMessageElement.text(GlobalResources.ERR_CO_LOI_XAY_RA);
            doiTuongDetail.append(errorMessageElement);

            ShowActionButtonDefault();
        },
        complete: function (xhr, status) {
            Common.UI.UnBlockElement("#detailPanel");
        }
    });
}

function ShowEditForm() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();
    var doiTuongDetail = $("#doiTuongDetail");
    Common.UI.BlockElement("#detailPanel");

    $.ajax({
        url: '/TiemChung/DoiTuong/Edit',
        type: "GET",
        async: true,
        data: { doiTuongId: doiTuongId },
        success: function (response) {
            doiTuongDetail.empty();
            doiTuongDetail.html(response);

            ShowActionButtonForEdit();
        },
        error: function (e) {
            doiTuongDetail.empty();

            // show error notification
            var errorMessageElement = $("<div />");
            errorMessageElement.attr("style", "margin:10px; text-align:center;");
            errorMessageElement.addClass("has-error");
            errorMessageElement.text(GlobalResources.ERR_CO_LOI_XAY_RA);
            doiTuongDetail.append(errorMessageElement);

            ShowActionButtonDefault();
        },
        complete: function (xhr, status) {
            Common.UI.UnBlockElement("#detailPanel");
        }
    });
}

function ConfirmDelete() {
    var doiTuongId = $("#hfDoiTuongId_Detail").val();

    // hien thi confirm xóa
    bootbox.confirm({
        size: 'small',
        message: GlobalResources.DT_Xoadoituongcanhbao,
        title: GlobalResources.DT_THONG_BAO,
        callback: function (result) {
            if (result != null) {
                if (result) {
                    AttemptDeleteDoiTuong(doiTuongId);
                }
            }
        },
        buttons: {
            confirm: {
                label: GlobalResources.DONG_Y,
                className: 'btn-hiden'
            },
            cancel: {
                label: GlobalResources.DONG,
                className: 'btn-default'
            }
        }
    });
}

function AttemptDeleteDoiTuong(doiTuongId) {
    Common.UI.BlockElement("#detailPanel");
    $.ajax({
        type: "POST",
        url: "/TiemChung/DoiTuong/Delete",
        data: { doiTuongId: doiTuongId },
        headers: layGiaTriToken(),
        async: true,
        success: function (response) {
            Common.UI.UnBlockElement("#detailPanel");
            if (response.Status == 0) {
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: response.Message,
                    class_name: 'growl-danger',
                    timeout: '',
                    sticky: false
                });
            } else if (response.Status == 1) {
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: response.Message,
                    class_name: 'growl-success',
                    timeout: '',
                    sticky: false
                });

                // refresh lai danh sach tim kiem
                OnShowPage(1, 20);

                var doiTuongDetail = $("#doiTuongDetail");
                doiTuongDetail.empty();

                var message = $("<div style='padding: 10px'>" + GlobalResources.DT_CHON_MOT_DOI_TUONG_DE_XEM_THONG_TIN + "</div>");
                doiTuongDetail.append(message);
                ShowActionButtonDefault();
            }
        },
        error: function (e) {
            Common.UI.UnBlockElement("#detailPanel");

            jQuery.gritter.add({
                title: GlobalResources.THONG_BAO,
                text: e.message,
                class_name: 'growl-danger',
                timeout: '',
                sticky: false
            });
        }
    });
}

/*** WARNING: Do not use this function. Please contact with HUNGBD2.***/
function SinhLaiMaDoiTuong() {
    Common.UI.BlockElement("#detailPanel");
    $.ajax({
        type: "POST",
        url: "/TiemChung/DoiTuong/SinhLaiMaDoiTuong",
        async: true,
        headers: layGiaTriToken(),
        success: function (response) {
            Common.UI.UnBlockElement("#detailPanel");
            if (response.Status == 0) {
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: response.Message,
                    class_name: 'growl-danger',
                    timeout: '',
                    sticky: false
                });
            } else if (response.Status == 1) {
                jQuery.gritter.add({
                    title: GlobalResources.THONG_BAO,
                    text: response.Message,
                    class_name: 'growl-success',
                    timeout: '',
                    sticky: false
                });
            }
        },
        error: function (e) {
            Common.UI.UnBlockElement("#detailPanel");

            jQuery.gritter.add({
                title: GlobalResources.THONG_BAO,
                text: e.message,
                class_name: 'growl-danger',
                timeout: '',
                sticky: false
            });
        }
    });
}

function PrintBarcode(MaDoiTuong, HoTen) {
    var baseUrl = '/TiemChung/DoiTuong/PrintBarcode';
    baseUrl += '?MA_DOI_TUONG=' + MaDoiTuong + '&HO_TEN=' + HoTen;
    $('#HO_TEN_IN_MA_VACH').val(HoTen);


    bootbox.dialog({
        title: GlobalResources.DT_IN_MA_VACH + ' - ' + HoTen + ' - ' + MaDoiTuong,
        message: '<div class="row"> \
                            <div class="rdio rdio-primary col-md-3"> \
                                <input id="rd4634" name="radio" value="4634" type="radio" onchange="OnPageSizeChange(\''+ MaDoiTuong + '\',' + '4634)">\
                                <label for="rd4634">'+ GlobalResources.DT_KHO_GIAY_4634_MM + '</label> \
                            </div> \
                            <div class="rdio rdio-primary col-md-3"> \
                                <input id="rd9234" name="radio" value="9234" type="radio" checked="checked" onchange="OnPageSizeChange(\'' + MaDoiTuong + '\',' + '9234)">\
                                <label for="rd9234">'+ GlobalResources.DT_KHO_GIAY_9234_MM + '</label> \
                            </div> \
                            <div class="rdio rdio-primary col-md-3"> \
                                <input id="rd300110" name="radio" value="300110" type="radio" onchange="OnPageSizeChange(\'' + MaDoiTuong + '\',' + '300110)">\
                                <label for="rd300110">' + GlobalResources.DT_KHO_GIAY_300110_MM + '</label> \
                            </div> \
                          </div> \
                          <div class="row mt5"><iframe id="frmPrintPreview" src="' + baseUrl + "&PAGE_SIZE=9234" + '" height="300" width="100%" border="0"></iframe></div>',
        size: 'large'
    });
}

function OnPageSizeChange(MaDoiTuong, size) {
    var baseUrl = '/TiemChung/DoiTuong/PrintBarcode';
    baseUrl += '?MA_DOI_TUONG=' + MaDoiTuong + '&HO_TEN=' + $('#HO_TEN_IN_MA_VACH').val();
    if ($("#rd4634").attr('checked')) {
        var completeUrl = baseUrl + "&PAGE_SIZE=" + 4634;
        $("#frmPrintPreview").attr("src", completeUrl);
    } else if ($("#rd9234").attr('checked')) {
        var completeUrl = baseUrl + "&PAGE_SIZE=" + 9234;
        $("#frmPrintPreview").attr("src", completeUrl);
    } else if ($("#rd300110").attr('checked')) {
        var completeUrl = baseUrl + "&PAGE_SIZE=" + 300110;
        $("#frmPrintPreview").attr("src", completeUrl);
    }
}

function Save() {
    $("#frmPhuNu").submit();
}

function Cancel() {
    var doiTuongDetail = $("#doiTuongDetail");
    doiTuongDetail.empty();

    var message = $("<div style='padding: 10px'>" + GlobalResources.DT_CHON_MOT_DOI_TUONG_DE_XEM_THONG_TIN + "</div>");
    doiTuongDetail.append(message);
    ShowActionButtonDefault();
}

function ShowActionButtonForDetail() {
    $("#btnAdd").show();
    $("#btnAddDTFromFile").show();
    $("#btnEdit").show();
    $("#btnDelete").show();
    $("#btnPrintBarcode").show();

    $("#btnSave").hide();
    $("#btnCancel").hide();

    $("#btnAdd2").show();
    $("#btnEdit2").show();
    $("#btnDelete2").show();
    $("#btnSave2").hide();
    $("#btnCancel2").hide();
}
function ShowActionButtonForDetailDoiTuongGop() {
    $("#btnAdd").show();
    $("#btnAddDTFromFile").show();
    $("#btnEdit").hide();
    $("#btnDelete").hide();
    $("#btnPrintBarcode").hide();

    $("#btnSave").hide();
    $("#btnCancel").hide();

    $("#btnAdd2").show();
    $("#btnEdit2").hide();
    $("#btnDelete2").hide();
    $("#btnSave2").hide();
    $("#btnCancel2").hide();

    $("#btnThemMuiTiem").hide();
    $("#btnThemNhanh").hide();

    $("#btnDangKyDichVu").hide();
    $('a[name="btnActionLST"]').hide();
}

function ShowActionButtonForAddNew() {
    $("#btnAdd").hide();
    $("#btnAddDTFromFile").hide();
    $("#btnEdit").hide();
    $("#btnDelete").hide();
    $("#btnPrintBarcode").hide();

    $("#btnSave").show();
    $("#btnCancel").show();

    $("#btnAdd2").hide();
    $("#btnEdit2").hide();
    $("#btnDelete2").hide();
    $("#btnSave2").show();
    $("#btnCancel2").show();
}

function ShowActionButtonForEdit() {
    $("#btnAdd").hide();
    $("#btnAddDTFromFile").hide();
    $("#btnEdit").hide();
    $("#btnDelete").hide();
    $("#btnPrintBarcode").hide();
    $("#btnSave").show();
    $("#btnCancel").show();

    $("#btnAdd2").hide();
    $("#btnEdit2").hide();
    $("#btnDelete2").hide();
    $("#btnSave2").show();
    $("#btnCancel2").show();
}

function ShowActionButtonDefault() {
    $("#btnAdd").show();
    $("#btnAddDTFromFile").show();
    $("#btnEdit").hide();
    $("#btnDelete").hide();
    $("#btnPrintBarcode").hide();
    $("#btnSave").hide();
    $("#btnCancel").hide();

    $("#btnAdd2").show();
    $("#btnEdit2").hide();
    $("#btnDelete2").hide();
    $("#btnSave2").hide();
    $("#btnCancel2").hide();
}

//----------- Search Form -------------//
function OnSearchButtonClick() {

    // reset page
    $("#hfSearchFormPageNumber").val(1);
    $("#hfSearchFormPageSize").val(20);
}

//----------- Search Result -----------//
function ShowNextPage() {
    // invoke callback from parent
    var divSearchResult = $("#divSearchResult");
    var pageNumber = $("#hfSearchResultPageNumber").val();
    var pageSize = $("#hfSearchResultPageSize").val();

    pageNumber++;
    OnShowPage(pageNumber, pageSize);
}

function ShowPreviousPage() {
    // invoke callback from parent
    var divSearchResult = $("#divSearchResult");
    var pageNumber = $("#hfSearchResultPageNumber").val();
    var pageSize = $("#hfSearchResultPageSize").val();

    pageNumber--;
    OnShowPage(pageNumber, pageSize);
}


//----------- Detail -----------//
function OnShowDetail(doiTuongId) {
    // xoa style cua doi tuong dang xem chi tiet
    $("#doiTuongSearchResult").find(".active-row").removeClass("active-row");

    // them style cho doi tuong moi
    $("#doiTuongSearchResult").find("[data-id='" + doiTuongId + "']").addClass("active-row");

    var doiTuongDetail = $("#doiTuongDetail");

    Common.UI.BlockElement("#detailPanel");

    $.ajax({
        url: '/TiemChung/DoiTuong/Detail',
        type: "GET",
        async: true,
        data: { doiTuongId: doiTuongId },
        success: function (response) {
            Common.UI.UnBlockElement("#detailPanel");

            doiTuongDetail.empty();
            doiTuongDetail.html(response);
            if (TiemChung$DoiTuong$Detail$DoiTuongModel != undefined && TiemChung$DoiTuong$Detail$DoiTuongModel != null && TiemChung$DoiTuong$Detail$DoiTuongModel.IS_ACTIVE == 1)
                ShowActionButtonForDetailDoiTuongGop();
            else
                ShowActionButtonForDetail();
        },
        error: function (e) {
            doiTuongDetail.empty();

            // show error notification
            var errorMessageElement = $("<div />");
            errorMessageElement.attr("style", "margin:10px; text-align:center;");
            errorMessageElement.addClass("has-error");
            errorMessageElement.text(GlobalResources.ERR_CO_LOI_XAY_RA);
            doiTuongDetail.append(errorMessageElement);

            ShowActionButtonDefault();
        },
        complete: function (xhr, status) {
            Common.UI.UnBlockElement("#detailPanel");
        }
    });
}

/// Migrate du lieu Audit ///
function MigrateDuLieu() {
    $.ajax({
        url: "/TiemChung/DoiTuong/MigrateDataAudit",
        type: "POST",
        async: true,
        success: function (response) {
            if (response.Status == 2) {
                jQuery.gritter.add({
                    text: GlobalResources.DT_MIGRATE_THANH_CONG,
                    class_name: 'growl-success',
                    sticky: false,
                    timeout: 2000
                });
            } else if (response.Status == 1) {
                var spResult = $('#spResult');
                var count = Number(spResult.text());
                count += response.Count;
                spResult.text(count);
            }
        },
        error: function (xhr, code, error) {
            jQuery.gritter.add({
                text: xhr.Message,
                class_name: 'growl-error',
                sticky: false,
                timeout: 2000
            });
        }
    });
}
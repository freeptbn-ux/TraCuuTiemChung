var ImCache = (function () {

    var exports = {};

    var Constant_Cache = {
        COSO_HINH_THUC: {
            TRUNG_UONG: 0,
            KHU_VUC: 1,
            TINH: 2,
            HUYEN: 3,
            XA: 4,
            CO_SO_DICH_VU: 5,
            BENH_VIEN: 7
        },
        COSO_PHAN_LOAI: {
            CO_SO_TIEM_CHUNG_CONG: 0,
            CO_SO_TIEM_CHUNG_TU: 1,
            BENH_VIEN: 2,
            VIETTEL_ICT: 8,
            KINH_DOANH_VTT: 9,
            KHAC: 3,
            KHONG_RO: -1
        }
    };

    //// VACXIN ////
    exports.GetDsVacxinFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("VACXIN_CACHED") == null) {
                // Caching danh muc Tinh, huyen xa
                $.ajax({
                    url: "/Vacxin/DsVacxin",
                    async: true,
                    type: "GET",
                    success: function (result) {

                        if (sessionStorage.getItem("VACXIN_CACHED") == null) {
                            sessionStorage.setItem("VACXIN_DATA", JSON.stringify(result));
                            sessionStorage.setItem("VACXIN_CACHED", "true");

                            callback(result);
                        }
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("VACXIN_DATA")));
            }
        } else {
            $.ajax({
                url: "/Vacxin/DsVacxin",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    exports.GetDsVacxinKhongCovidFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("VACXINKHONGCOVID_CACHED") == null) {
                // Caching danh muc Tinh, huyen xa
                $.ajax({
                    url: "/Vacxin/DsVacxinKhongCovid",
                    async: true,
                    type: "GET",
                    success: function (result) {

                        if (sessionStorage.getItem("VACXINKHONGCOVID_CACHED") == null) {
                            sessionStorage.setItem("VACXINKHONGCOVID_DATA", JSON.stringify(result));
                            sessionStorage.setItem("VACXINKHONGCOVID_CACHED", "true");

                            callback(result);
                        }
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("VACXINKHONGCOVID_DATA")));
            }
        } else {
            $.ajax({
                url: "/Vacxin/DsVacxinKhongCovid",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    exports.GetDsVacxinCovidFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("VACXINCOVID_CACHED") == null) {
                // Caching danh muc Tinh, huyen xa
                $.ajax({
                    url: "/Vacxin/DsVacxinCovid",
                    async: true,
                    type: "GET",
                    success: function (result) {

                        if (sessionStorage.getItem("VACXINCOVID_CACHED") == null) {
                            sessionStorage.setItem("VACXINCOVID_DATA", JSON.stringify(result));
                            sessionStorage.setItem("VACXINCOVID_CACHED", "true");

                            callback(result);
                        }
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("VACXINCOVID_DATA")));
            }
        } else {
            $.ajax({
                url: "/Vacxin/DsVacxinCovid",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    exports.GetVacxinByIdFromCache = function (vacxinId, callback) {
        exports.GetDsVacxinFromCache(function (dsVacxin) {
            var vacxinItem = dsVacxin.filter(function (v) {
                if (v.VACXIN_ID == vacxinId) {
                    return v;
                }
            });

            if (vacxinItem != null && vacxinItem.length > 0) {
                callback(vacxinItem[0]);
            } else {
                callback(null);
            }
        });
    }
    //// QUOC GIA ////
    exports.GetDsQuocGiaFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("QUOC_GIA_CACHED") == null) {
                // Caching danh muc Dân tộc
                $.ajax({
                    url: "/DanToc/DsQuocGia",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("QUOC_GIA_DATA", JSON.stringify(result));
                        sessionStorage.setItem("QUOC_GIA_CACHED", "true");

                        callback(result);
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("QUOC_GIA_DATA")));
            }
        } else {
            $.ajax({
                url: "/DanToc/DsQuocGia",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    //// DAN TOC ////
    exports.GetDsDanTocFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("DAN_TOC_CACHED") == null) {
                // Caching danh muc Dân tộc
                $.ajax({
                    url: "/DanToc/DsDanToc",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("DAN_TOC_DATA", JSON.stringify(result));
                        sessionStorage.setItem("DAN_TOC_CACHED", "true");

                        callback(result);
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("DAN_TOC_DATA")));
            }
        } else {
            $.ajax({
                url: "/DanToc/DsDanToc",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    exports.GetDanTocByIdFromCache = function (danTocId, callback) {
        exports.GetDsDanTocFromCache(function (dsDanToc) {
            var danToc = dsDanToc.filter(function (danTocItem) {
                if (danTocItem.DAN_TOC_ID == danTocId) {
                    return danTocItem;
                }
            });

            if (danToc != null && danToc.length > 0) {
                callback(danToc[0]);
            } else {
                callback(null);
            }
        });
    }

    //// TÔN GIÁO ////
    exports.GetDsTonGiaoFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            // Kiểm tra xem có cache chưa
            if (sessionStorage.getItem("TON_GIAO_CACHED") == null) {
                // Nếu chưa có thì gọi API lấy danh sách Tôn giáo
                $.ajax({
                    url: "/TonGiao/DsTonGiao",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        // Lưu dữ liệu vào sessionStorage
                        sessionStorage.setItem("TON_GIAO_DATA", JSON.stringify(result));
                        sessionStorage.setItem("TON_GIAO_CACHED", "true");

                        // Gọi callback trả về dữ liệu
                        callback(result);
                    },
                    error: function (xhr, status, error) {
                        console.error("Lỗi khi tải danh mục Tôn giáo:", error);
                        callback([]);
                    }
                });
            } else {
                // Nếu có cache thì lấy từ sessionStorage
                callback(JSON.parse(sessionStorage.getItem("TON_GIAO_DATA")));
            }
        } else {
            // Trường hợp trình duyệt không hỗ trợ sessionStorage
            $.ajax({
                url: "/TonGiao/DsTonGiao",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                },
                error: function (xhr, status, error) {
                    console.error("Lỗi khi tải danh mục Tôn giáo:", error);
                    callback([]);
                }
            });
        }
    };

    // Hàm lấy chi tiết Tôn giáo theo ID từ cache
    exports.GetTonGiaoByIdFromCache = function (tonGiaoId, callback) {
        exports.GetDsTonGiaoFromCache(function (dsTonGiao) {
            var tonGiao = dsTonGiao?.filter(function (tonGiaoItem) {
                return tonGiaoItem.TON_GIAO_ID == tonGiaoId;
            });

            if (tonGiao != null && tonGiao.length > 0) {
                callback(tonGiao[0]);
            } else {
                callback(null);
            }
        });
    };

    //// VUNG MIEN ////
    exports.GetDsVungMienFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("VUNG_MIEN_CACHED") == null) {
                // Caching danh muc Vung Mien
                $.ajax({
                    url: "/DonViHanhChinh/DsVungMien",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("VUNG_MIEN_DATA", JSON.stringify(result));
                        sessionStorage.setItem("VUNG_MIEN_CACHED", "true");

                        callback(result.sort(function (a, b) {
                            return (a.TEN_VUNG_MIEN > b.TEN_VUNG_MIEN);
                        }));
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("VUNG_MIEN_DATA")).sort(function (a, b) {
                    return (a.TEN_VUNG_MIEN > b.TEN_VUNG_MIEN);
                }));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsVungMien",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TEN_VUNG_MIEN > b.TEN_VUNG_MIEN);
                    }));
                }
            });
        }
    }

    //// DM DON VI HANH CHINH////
    exports.GetDsDonViHanhChinhFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("DMDONVIHANHCHINH_CACHED") == null) {
                // Caching danh muc don vi hanh chinh
                $.ajax({
                    url: "/DonViHanhChinh/DsDonViHanhChinh",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("DMDONVIHANHCHINH_DATA", JSON.stringify(result));
                        sessionStorage.setItem("DMDONVIHANHCHINH_CACHED", "true");

                        callback(result.sort(function (a, b) {
                            return (a.TENDAYDU > b.TENDAYDU);
                        }));
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("DMDONVIHANHCHINH_DATA")).sort(function (a, b) {
                    return (a.TENDAYDU > b.TENDAYDU);
                }));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsDonViHanhChinh",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TENDAYDU > b.TENDAYDU);
                    }));
                }
            });
        }
    }

    exports.GetDsDonViHanhChinhByVungMienFromCache = function (vungMienId, callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("DMDONVIHANHCHINH_CACHED") == null || sessionStorage.getItem("DMDONVIHANHCHINH_CACHED") == undefined) {
                $.ajax({
                    url: "/DonViHanhChinh/DsDonViHanhChinhByVungMienId",
                    async: true,
                    data: { vungMienId: -1 },
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("DMDONVIHANHCHINH_DATA", JSON.stringify(result));
                        sessionStorage.setItem("DMDONVIHANHCHINH_CACHED", "true");

                        if (vungMienId == -1) {
                            callback(result.sort(function (a, b) {
                                return (a.TENDAYDU > b.TENDAYDU);
                            }));
                        } else {
                            callback(result.filter(function (dvhcItem) {
                                if (dvhcItem.NIISID == vungMienId) {
                                    return dvhcItem;
                                }
                            }).sort(function (a, b) {
                                return (a.TENDAYDU > b.TENDAYDU);
                            }));
                        }
                    }
                });
            } else {
                var dsTinhCached = JSON.parse(sessionStorage.getItem("DMDONVIHANHCHINH_DATA"));
                if (vungMienId == -1) {
                    callback(dsTinhCached.sort(function (a, b) {
                        return (a.TENDAYDU > b.TENDAYDU);
                    }));
                } else {
                    callback(dsTinhCached.filter(function (dvhcItem) {
                        if (dvhcItem.NIISID == vungMienId) {
                            return dvhcItem;
                        }
                    }).sort(function (a, b) {
                        return (a.TENDAYDU > b.TENDAYDU);
                    }));
                }
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsDonViHanhChinhByVungMienId",
                async: true,
                type: "GET",
                data: { vungMienId: vungMienId },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TENDAYDU > b.TENDAYDU);
                    }));
                }
            });
        }
    }

    //// TINH ////
    exports.GetDsTinhFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("TINH_CACHED") == null) {
                // Caching danh muc Tinh
                $.ajax({
                    url: "/DonViHanhChinh/DsTinh",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("TINH_DATA", JSON.stringify(result));
                        sessionStorage.setItem("TINH_CACHED", "true");

                        callback(result.sort(function (a, b) {
                            return (a.TENTINH > b.TENTINH);
                        }));
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("TINH_DATA")).sort(function (a, b) {
                    return (a.TENTINH > b.TENTINH);
                }));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsTinh",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TENTINH > b.TENTINH);
                    }));
                }
            });
        }
    }

    exports.GetDsTinhByVungMienFromCache = function (vungMienId, callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("TINH_CACHED") == null) {
                // Caching danh muc Tỉnh
                $.ajax({
                    url: "/DonViHanhChinh/DsTinhByVungMienId",
                    async: true,
                    data: { vungMienId: -1 },
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("TINH_DATA", JSON.stringify(result));
                        sessionStorage.setItem("TINH_CACHED", "true");

                        if (vungMienId == -1) {
                            callback(result.sort(function (a, b) {
                                return (a.TENTINH > b.TENTINH);
                            }));
                        } else {
                            callback(result.filter(function (tinhItem) {
                                if (tinhItem.VUNG_MIEN == vungMienId) {
                                    return tinhItem;
                                }
                            }).sort(function (a, b) {
                                return (a.TENTINH > b.TENTINH);
                            }));
                        }
                    }
                });
            } else {
                var dsTinhCached = JSON.parse(sessionStorage.getItem("TINH_DATA"));
                if (vungMienId == -1) {
                    callback(dsTinhCached.sort(function (a, b) {
                        return (a.TENTINH > b.TENTINH);
                    }));
                } else {
                    callback(dsTinhCached.filter(function (tinhItem) {
                        if (tinhItem.VUNG_MIEN == vungMienId) {
                            return tinhItem;
                        }
                    }).sort(function (a, b) {
                        return (a.TENTINH > b.TENTINH);
                    }));
                }
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsTinhByVungMienId",
                async: true,
                type: "GET",
                data: { vungMienId: vungMienId },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TENTINH > b.TENTINH);
                    }));
                }
            });
        }
    }

    exports.GetTinhByIdFromCache = function (tinhId, callback) {
        //console.log("TINH ID" + tinhId);
        exports.GetDsTinhByVungMienFromCache(-1, function (dsTinh) {
            var tinh = dsTinh.filter(function (tinhItem) {
                if (tinhItem.TINH_ID == tinhId) {
                    return tinhItem;
                }
            });

            if (tinh != null && tinh.length > 0) {
                callback(tinh[0]);
            } else {
                callback(null);
            }
        });
    }

    //// HUYEN ////
    exports.GetDsHuyenFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("HUYEN_CACHED") == null) {
                // Caching danh muc Huyen
                $.ajax({
                    url: "/DonViHanhChinh/DsHuyen",
                    async: true,
                    type: "GET",
                    data: { tinhId: -1 },
                    success: function (result) {
                        sessionStorage.setItem("HUYEN_DATA", JSON.stringify(result));
                        sessionStorage.setItem("HUYEN_CACHED", "true");

                        callback(result.sort(function (a, b) {
                            return (a.TENHUYEN > b.TENHUYEN);
                        }));
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("HUYEN_DATA")).sort(function (a, b) {
                    return (a.TENHUYEN > b.TENHUYEN);
                }));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsHuyen",
                async: true,
                type: "GET",
                data: { tinhId: -1 },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TENHUYEN > b.TENHUYEN);
                    }));
                }
            });
        }
    }

    exports.GetDsHuyenByTinhFromCache = function (tinhId, callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("HUYEN_CACHED") == null) {
                // Caching danh muc Huyện
                $.ajax({
                    url: "/DonViHanhChinh/DsHuyen",
                    async: true,
                    data: { tinhId: -1 },
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("HUYEN_DATA", JSON.stringify(result));
                        sessionStorage.setItem("HUYEN_CACHED", "true");

                        if (tinhId == -1) {
                            callback(result.sort(function (a, b) {
                                return (a.TENHUYEN > b.TENHUYEN);
                            }));
                        } else {
                            callback(result.filter(function (huyenItem) {
                                if (huyenItem.TINH_ID == tinhId) {
                                    return huyenItem;
                                }
                            }).sort(function (a, b) {
                                return (a.TENHUYEN > b.TENHUYEN);
                            }));
                        }
                    }
                });
            } else {
                var dsHuyenCached = JSON.parse(sessionStorage.getItem("HUYEN_DATA"));
                if (tinhId == -1) {
                    callback(dsHuyenCached.sort(function (a, b) {
                        return (a.TENHUYEN > b.TENHUYEN);
                    }));
                } else {
                    callback(dsHuyenCached.filter(function (huyenItem) {
                        if (huyenItem.TINH_ID == tinhId) {
                            return huyenItem;
                        }
                    }).sort(function (a, b) {
                        return (a.TENHUYEN > b.TENHUYEN);
                    }));
                }
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsHuyen",
                async: true,
                type: "GET",
                data: { tinhId: tinhId },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TENHUYEN > b.TENHUYEN);
                    }));
                }
            });
        }
    }

    exports.GetHuyenByIdFromCache = function (huyenId, callback) {
        exports.GetDsHuyenByTinhFromCache(-1, function (dsHuyen) {
            var huyen = dsHuyen.filter(function (huyenItem) {
                if (huyenItem.HUYEN_ID == huyenId) {
                    return huyenItem;
                }
            });

            if (huyen != null && huyen.length > 0) {
                callback(huyen[0]);
            } else {
                callback(null);
            }
        });
    }

    //// XA ////
    exports.GetDsXaFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("XA_CACHED") == null) {
                // Caching danh muc Xa
                $.ajax({
                    url: "/DonViHanhChinh/DsXa",
                    async: true,
                    type: "GET",
                    data: { huyenId: -1 },
                    success: function (result) {
                        sessionStorage.setItem("XA_DATA", JSON.stringify(result));
                        sessionStorage.setItem("XA_CACHED", "true");

                        callback(result.sort(function (a, b) {
                            return (a.TEN_XA > b.TEN_XA);
                        }));
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("XA_DATA")).sort(function (a, b) {
                    return (a.TEN_XA > b.TEN_XA);
                }));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsXa",
                async: true,
                type: "GET",
                data: { huyenId: -1 },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                }
            });
        }
    }

    exports.GetDsXaByHuyenFromCache = function (huyenId, callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("XA_CACHED") == null) {
                // Caching danh muc Xa
                $.ajax({
                    url: "/DonViHanhChinh/DsXa",
                    async: true,
                    data: { huyenId: -1 },
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("XA_DATA", JSON.stringify(result));
                        sessionStorage.setItem("XA_CACHED", "true");

                        if (huyenId == -1) {
                            callback(result.sort(function (a, b) {
                                return (a.TEN_XA > b.TEN_XA);
                            }));
                        } else {
                            callback(result.filter(function (xaItem) {
                                if (xaItem.HUYEN_ID == huyenId) {
                                    return xaItem;
                                }
                            }).sort(function (a, b) {
                                return (a.TEN_XA > b.TEN_XA);
                            }));
                        }
                    }
                });
            } else {
                var dsXaCached = JSON.parse(sessionStorage.getItem("XA_DATA"));
                if (huyenId == -1) {
                    callback(dsXaCached.sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                } else {
                    callback(dsXaCached.filter(function (xaItem) {
                        if (xaItem.HUYEN_ID == huyenId) {
                            return xaItem;
                        }
                    }).sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                }
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsXa",
                async: true,
                type: "GET",
                data: { huyenId: huyenId },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                }
            });
        }
    }

    exports.GetDsXa715FromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("XA_CACHED") == null) {
                // Caching danh muc Xa
                $.ajax({
                    url: "/DonViHanhChinh/DsXa715",
                    async: true,
                    type: "GET",
                    data: { tinhId: -1 },
                    success: function (result) {
                        sessionStorage.setItem("XA_DATA", JSON.stringify(result));
                        sessionStorage.setItem("XA_CACHED", "true");

                        callback(result.sort(function (a, b) {
                            return (a.TEN_XA > b.TEN_XA);
                        }));
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("XA_DATA")).sort(function (a, b) {
                    return (a.TEN_XA > b.TEN_XA);
                }));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsXa715",
                async: true,
                type: "GET",
                data: { tinhId: -1 },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                }
            });
        }
    }

    exports.GetDsXaByTinhFromCache = function (tinhId, callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("XA_CACHED") == null) {
                // Caching danh muc Xa
                $.ajax({
                    url: "/DonViHanhChinh/DsXa715",
                    async: true,
                    data: { tinhId: -1 },
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("XA_DATA", JSON.stringify(result));
                        sessionStorage.setItem("XA_CACHED", "true");

                        if (tinhId == -1) {
                            callback(result.sort(function (a, b) {
                                return (a.TEN_XA > b.TEN_XA);
                            }));
                        } else {
                            callback(result.filter(function (xaItem) {
                                if (xaItem.TINH_ID == tinhId) {
                                    return xaItem;
                                }
                            }).sort(function (a, b) {
                                return (a.TEN_XA > b.TEN_XA);
                            }));
                        }
                    }
                });
            } else {
                var dsXaCached = JSON.parse(sessionStorage.getItem("XA_DATA"));
                if (tinhId == -1) {
                    callback(dsXaCached.sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                } else {
                    callback(dsXaCached.filter(function (xaItem) {
                        if (xaItem.TINH_ID == tinhId) {
                            return xaItem;
                        }
                    }).sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                }
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsXa715",
                async: true,
                type: "GET",
                data: { tinhId: tinhId },
                success: function (result) {
                    callback(result.sort(function (a, b) {
                        return (a.TEN_XA > b.TEN_XA);
                    }));
                }
            });
        }
    }

    exports.GetXaByIdFromCache = function (xaId, callback) {
        if (!xaId) {
            callback(null);
            return;
        }

        exports.GetDsXaByHuyenFromCache(-1, function (dsXa) {
            if (!Array.isArray(dsXa) || dsXa.length === 0) {
                callback(null);
                return;
            }

            var xa = dsXa.find(function (item) {
                return item && item.XA_ID == xaId;
            });

            callback(xa || null);
        });
    }

    //// CO SO TIEM CHUNG ////
    exports.GetDsCoSoTiemChungFromCache = function (tinhId, callback) {
        console.log('tinhId: ', tinhId);
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("CO_SO_TIEM_CHUNG_CACHED") == null) {
                // Caching danh muc cơ sở tiêm chủng
                $.ajax({
                    url: "/DonViHanhChinh/DsCoSoTiemChung",
                    async: true,
                    data: { tinhId },
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("CO_SO_TIEM_CHUNG_DATA", JSON.stringify(result));
                        sessionStorage.setItem("CO_SO_TIEM_CHUNG_CACHED", "true");

                        callback(result);
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("CO_SO_TIEM_CHUNG_DATA")));
            }
        } else {
            console.log('tinhId: ', tinhId);
            $.ajax({
                url: "/DonViHanhChinh/DsCoSoTiemChung",
                async: true,
                data: { tinhId },
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    exports.GetCoSoTiemChungByIdFromCache = function (coSoId, callback) {
        exports.GetDsCoSoTiemChungFromCache(function (dsCoso) {
            var coSo = dsCoso.filter(function (coSoItem) {
                if (coSoItem.COSO_ID == coSoId) {
                    return coSoItem;
                }
            });

            if (coSo != null && coSo.length > 0) {
                callback(coSo[0]);
            } else {
                callback(null);
            }
        });
    }

    exports.GetDsCoSoTiemChungByParamsFromCache = function (tinhId, huyenId, xaId, hinhThuc, vungMien, callback) {
        exports.GetDsCoSoTiemChungFromCache(tinhId, function (listObjs) {
            var lstResult;
            if (hinhThuc == Constant_Cache.COSO_HINH_THUC.TRUNG_UONG) // TW
            {
                lstResult = listObjs.filter(o => o.HINHTHUC == Constant_Cache.COSO_HINH_THUC.TRUNG_UONG);
            }
            else if (hinhThuc == Constant_Cache.COSO_HINH_THUC.KHU_VUC) // Khu Vuc
            {
                lstResult = listObjs.filter(o => o.HINHTHUC == Constant_Cache.COSO_HINH_THUC.KHU_VUC && o.VUNG_MIEN_ID == vungMien);
            }
            else if (hinhThuc == Constant_Cache.COSO_HINH_THUC.TINH) // Tinh
            {
                lstResult = listObjs.filter(o => {
                    if (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_CONG && o.HINHTHUC == Constant_Cache.COSO_HINH_THUC.TINH && o.TINH_ID == tinhId) {
                        return o;
                    }
                });
            }
            else if (hinhThuc == Constant_Cache.COSO_HINH_THUC.HUYEN) // Huyen
            {
                lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_CONG) && (o.HINHTHUC == Constant_Cache.COSO_HINH_THUC.HUYEN) && (o.HUYEN_ID == huyenId));
            }
            else if (hinhThuc == Constant_Cache.COSO_HINH_THUC.XA) // Xa
            {
                lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_CONG) && (o.HINHTHUC == Constant_Cache.COSO_HINH_THUC.XA) && (o.XA_ID == xaId));
            }
            else if (hinhThuc == Constant_Cache.COSO_HINH_THUC.CO_SO_DICH_VU)//trường hợp cstc dịch vụ
            {
                if (xaId > 0) {
                    lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_TU) && (o.TINH_ID == tinhId) && (o.HUYEN_ID == huyenId) && (o.XA_ID == xaId));
                }
                else if (huyenId > 0) {
                    lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_TU) && (o.TINH_ID == tinhId) && (o.HUYEN_ID == huyenId));
                }
                else {
                    lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_TU) && (o.TINH_ID == tinhId));
                }

            }
            else if (hinhThuc == Constant_Cache.COSO_HINH_THUC.BENH_VIEN)//trường hợp benh vien
            {
                if (xaId > 0) {
                    lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.BENH_VIEN) && (o.TINH_ID == tinhId) && (o.HUYEN_ID == huyenId) && (o.XA_ID == xaId));
                }
                else if (huyenId > 0) {
                    lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.BENH_VIEN) && (o.TINH_ID == tinhId) && (o.HUYEN_ID == huyenId));
                }
                else {
                    lstResult = listObjs.filter(o => (o.PHANLOAI == Constant_Cache.COSO_PHAN_LOAI.BENH_VIEN) && (o.TINH_ID == tinhId));
                }
            }

            if (lstResult && lstResult.length > 0) {
                callback(lstResult);
            } else {
                callback(lstResult);
            }

        });
    }

    exports.GetDsCoSoTiemChungByTinhHuyenXaFromCache = function (tinhId, huyenId, xaId, callback) {
        exports.GetDsCoSoTiemChungFromCache(tinhId, function (listObjs) {
            var lstResult;
            if (xaId > 0) {
                lstResult = listObjs.filter(o => (o.TINH_ID == tinhId) && (o.HUYEN_ID == huyenId) && (o.XA_ID == xaId));
            }
            else if (huyenId > 0) {
                lstResult = listObjs.filter(o => (o.TINH_ID == tinhId) && (o.HUYEN_ID == huyenId));
            }
            else {
                lstResult = listObjs.filter(o => (o.TINH_ID == tinhId));
            }

            if (lstResult && lstResult.length > 0) {
                callback(lstResult);
            } else {
                callback(lstResult);
            }

        });
    }

    //// QUOC GIA ////
    exports.GetDsQuocGiaFromCache = function (callback) {
        if (typeof (Storage) !== "undefined") {
            if (sessionStorage.getItem("QUOC_GIA_CACHED") == null) {
                // Caching các quốc gia
                $.ajax({
                    url: "/DonViHanhChinh/DsQuocGia",
                    async: true,
                    type: "GET",
                    success: function (result) {
                        sessionStorage.setItem("QUOC_GIA_DATA", JSON.stringify(result));
                        sessionStorage.setItem("QUOC_GIA_CACHED", "true");

                        callback(result);
                    }
                });
            } else {
                callback(JSON.parse(sessionStorage.getItem("QUOC_GIA_DATA")));
            }
        } else {
            $.ajax({
                url: "/DonViHanhChinh/DsQuocGia",
                async: true,
                type: "GET",
                success: function (result) {
                    callback(result);
                }
            });
        }
    }

    exports.GetQuocGiaByIdFromCache = function (quocGiaId, callback) {
        exports.GetDsQuocGiaFromCache(function (dsQuocGia) {
            var quocGia = dsQuocGia.filter(function (quocGiaItem) {
                if (quocGiaItem.QUOC_GIA_ID == quocGiaId) {
                    return quocGiaItem;
                }
            });

            if (quocGia != null && quocGia.length > 0) {
                callback(quocGia[0]);
            } else {
                callback(null);
            }
        });
    }

    return exports;
})();
import React from "react";

import { Image, Spin } from "antd";
import { RcFile } from "antd/es/upload";
import { FormInstance } from "antd/lib";
import { useTranslation } from "next-i18next";
import { useRouter } from "next/router";

import FormCard from "@/components/common/form/formCard";
import DirectUploadCustom from "@/components/common/upload/DirectUploadCustom";
import DropZoneUpload from "@/components/common/upload/DropZoneUpload";
import { PathName } from "@/components/layouts/routeMapMember";
import useFileUpload from "@/hooks/common/useFileUpload";
import { IBank } from "@/models/bank/banks";
import { IImageLinkProps } from "@/models/product/product";
import { formatImageURL } from "@/utils/formatString";
import getBase64 from "@/utils/getbase64";

interface Props {
  bankDetail?: IBank;
  imageLink?: string;
  setImageLink: (value?: string) => void;
  isLoading: boolean;
  form: FormInstance;
}
const BankUploadImageCard = ({
  isLoading,
  bankDetail,
  imageLink,
  setImageLink,
  form,
}: Props) => {
  const { t } = useTranslation("bank");

  const { query } = useRouter();
  const bankId = query?.banksId?.toString() || "";
  const isCreate = bankId === PathName.new;

  const { onUploadFiles } = useFileUpload();
  // const isSelectedImage = Boolean(bankDetail?.logo && imageLink?.imageLink);

  // const imageUpdateSrc = isSelectedImage ? imageLink?.imageLinkValue : bankDetail?.logo && formatImageURL(bankDetail?.logo);
  // const showImagePreview = form?.getFieldValue('logo');

  return (
    <>
      <Spin spinning={isLoading}>
        <DropZoneUpload
          singleUpload
          currentImage={form.getFieldValue("logo")}
          form={form}
          detailImage={"https://fncdn.b-cdn.net/test/apple.jpg"}
          imageLink={imageLink}
        />
        {/* <FormCard
          titleText={t('image')}
          extra={
            <DirectUploadCustom
              hasImage={!!showImagePreview}
              onUpdate={async (imageFormValue) => {
                const uploadedImage = await onUploadFiles(imageFormValue?.fileList);

                const imageLinkValue = (await getBase64(imageFormValue?.fileList?.[0]?.originFileObj as RcFile)) as string;

                if (uploadedImage?.ok) {
                  setImageLink({ imageLinkValue, name: imageFormValue?.fileList?.[0]?.name, imageLink: uploadedImage?.attachments?.[0] });
                  form.setFieldValue('logo', uploadedImage?.attachments?.[0] || undefined);
                }
              }}
              onRemove={() => {
                form.setFieldValue('logo', undefined);
                setImageLink({
                  imageLinkValue: '',
                  name: '',
                  imageLink: '',
                });
              }}
              isNoPrompt
            />
          }
        >
          {showImagePreview ? (
            <Image.PreviewGroup>
              <Image src={!isCreate && isSelectedImage ? imageUpdateSrc : imageLink?.imageLinkValue} alt="image" />
            </Image.PreviewGroup>
          ) : (
            <div>{t(`no_image_selected`)}</div>
          )}
        </FormCard> */}
      </Spin>
    </>
  );
};

export default BankUploadImageCard;

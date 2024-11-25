import React from 'react';
import TransactionForm from '@/components/transaction';
import { ContentLayout } from "@/components/admin-panel/content-layout";
export default  function TransactionPage({ params: { lng,slug } }) {
    return (
            <ContentLayout title="Transaction">
            <TransactionForm lng={lng} slug={slug} />
         </ContentLayout>
    );
};


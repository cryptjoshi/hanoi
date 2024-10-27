import React from 'react';
import TransactionForm from '@/components/deposit';
import { ContentLayout } from "@/components/admin-panel/content-layout";
export default  function TransactionPage({ params: { lng } }) {
    return (
            <ContentLayout title="Transaction">
            <TransactionForm lng={lng} />
         </ContentLayout>
    );
};

